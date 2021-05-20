package agent

import (
	"context"
	"sync/atomic"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/log"
	"go.etcd.io/etcd/clientv3"
)

// Agent watches etcd database and synchornizes system
// users and groups with the database entries.
type Agent struct {
	*clientv3.Client
	etcdpasswd.Syncer
	rev int64
}

// StartWatching is a goroutine to watch etcd and notify updater.
func (a *Agent) StartWatching(ctx context.Context, updateCh chan<- struct{}) error {
	// obtain the current revision to avoid missing events
	// it's OK that key "/" does not exists
	resp, err := a.Get(ctx, "/")
	if err != nil {
		return err
	}

	// notify updater for the initial sync
	rev := resp.Header.Revision
	atomic.StoreInt64(&a.rev, rev)
	updateCh <- struct{}{}

	rch := a.Watch(ctx, "",
		clientv3.WithPrefix(),
		clientv3.WithRev(rev+1),
		clientv3.WithProgressNotify(),
	)
	for wresp := range rch {
		atomic.StoreInt64(&a.rev, wresp.Header.Revision)

		// notify updater if possible
		select {
		case updateCh <- struct{}{}:
		default:
			// Do nothing if previous notification has not been processed,
			// because updater reflects all updates by one notification.
		}
	}

	return nil
}

// StartUpdater is a goroutine to receive notification from watcher.
func (a *Agent) StartUpdater(ctx context.Context, updateCh <-chan struct{}) error {
	for {
		select {
		case <-updateCh:
			rev := atomic.LoadInt64(&a.rev)
			log.Info("start sync", map[string]interface{}{
				"rev": rev,
			})

			db, err := etcdpasswd.GetDatabase(ctx, a.Client, rev)
			if err != nil {
				return err
			}
			err = synchronize(ctx, db, a.Syncer)
			if err != nil {
				return err
			}

			log.Info("finish sync", map[string]interface{}{
				"rev": rev,
			})
		case <-ctx.Done():
			return nil
		}
	}
}

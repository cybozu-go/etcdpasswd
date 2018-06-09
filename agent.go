package etcdpasswd

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/cybozu-go/log"
)

// Agent watches etcd database and synchornizes system
// users and groups with the database entries.
type Agent struct {
	*clientv3.Client
	Syncer
}

// StartWatching is a goroutine to watch etcd and notify updater.
func (a Agent) StartWatching(ctx context.Context, updateCh chan<- int64) error {
	// obtain the current revision to avoid missing events
	// it's OK that key "/" does not exists
	resp, err := a.Get(ctx, "/")
	if err != nil {
		return err
	}

	rev := resp.Header.Revision

	// notify updater for the initial sync
	updateCh <- rev

	rch := a.Watch(ctx, "",
		clientv3.WithPrefix(),
		clientv3.WithRev(rev+1),
		clientv3.WithProgressNotify(),
	)
	for wresp := range rch {
		// notify updater if possible
		select {
		case updateCh <- wresp.Header.Revision:
		default:
			// Do nothing if previous notification has not been processed,
			// because updater reflects all updates by one notification.
		}
	}

	return nil
}

// StartUpdater is a goroutine to receive notification from watcher.
func (a Agent) StartUpdater(ctx context.Context, updateCh <-chan int64) error {
	for {
		select {
		case rev := <-updateCh:
			log.Info("start sync", map[string]interface{}{
				"rev": rev,
			})
			err := a.sync(ctx, rev)
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

func (a Agent) sync(ctx context.Context, rev int64) error {
	db, err := GetDatabase(ctx, a.Client, rev)
	if err != nil {
		return err
	}

	return db.Sync(ctx, a.Syncer)
}

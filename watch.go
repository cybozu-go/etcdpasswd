package etcdpasswd

import (
	"context"
	"encoding/json"
	"strings"
	"sync/atomic"

	"github.com/coreos/etcd/clientv3"
)

var config atomic.Value

func (c Client) handleConfig(ev *clientv3.Event) error {
	if ev.Type == clientv3.EventTypeDelete {
		return nil
	}

	cfg := new(Config)
	err := json.Unmarshal(ev.Kv.Value, cfg)
	if err != nil {
		return err
	}

	config.Store(cfg)
	return nil
}

// StartWatching starts watching for user/group modification on etcd database.
func (c Client) StartWatching(ctx context.Context, updateCh chan<- struct{}) error {
	// obtain the current revision to avoid missing events
	// it's OK that key "/" does not exists
	resp, err := c.Get(ctx, "/")
	if err != nil {
		return err
	}
	rev := resp.Header.Revision

	cfg, _, err := c.GetConfig(ctx)
	if err != nil {
		return err
	}
	config.Store(cfg)

	rch := c.Watch(ctx, "",
		clientv3.WithPrefix(),
		clientv3.WithRev(rev),
		clientv3.WithProgressNotify(),
	)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			var err error
			key := string(ev.Kv.Key)
			switch {
			case strings.HasPrefix(key, KeyLastUID) || strings.HasPrefix(key, KeyLastGID):
				// do nothing
			case key == KeyConfig:
				err = c.handleConfig(ev)
			default:
				// notify updater if possible
				select {
				case updateCh <- struct{}{}:
				default:
					// Do nothing if previous notification has not been processed,
					// because updater reflects all updates by one notification.
				}
			}
			if err != nil {
				panic(err)
			}
		}
		if len(wresp.Events) == 0 {
			// periodic notification
			select {
			case updateCh <- struct{}{}:
			default:
			}
		}
	}

	return nil
}

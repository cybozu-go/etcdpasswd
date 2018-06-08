package etcdpasswd

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/cybozu-go/log"
)

func handleDeletedUsersResp(ctx context.Context, resp *etcdserverpb.RangeResponse) error {
	for _, kv := range resp.Kvs {
		username := string(kv.Key)[len(KeyDeletedUsers+"/"):]
		err := auth.DeleteUser(ctx, username)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c Client) updateAll(ctx context.Context) error {
	log.Info("updating users and groups...", nil)

	resp, err := c.Txn(ctx).If().Then(
		clientv3.OpGet(KeyUsers, clientv3.WithPrefix()),
		clientv3.OpGet(KeyDeletedUsers, clientv3.WithPrefix()),
		clientv3.OpGet(KeyGroups, clientv3.WithPrefix()),
		clientv3.OpGet(KeyDeletedGroups, clientv3.WithPrefix()),
		clientv3.OpGet(KeyLocked, clientv3.WithPrefix()),
	).Commit()
	if err != nil {
		return err
	}

	//usersResp := resp.Responses[0].GetResponseRange()
	deletedUsersResp := resp.Responses[1].GetResponseRange()
	// groupsResp := resp.Responses[2].GetResponseRange()
	// deletedGroupsResp := resp.Responses[3].GetResponseRange()
	// lockedResp := resp.Responses[4].GetResponseRange()

	// delete users before deleting groups
	err = handleDeletedUsersResp(ctx, deletedUsersResp)
	if err != nil {
		return err
	}
	// err = handleDeletedGroupsResp(ctx, deletedGroupsResp)
	// if err != nil {
	// 	return err
	// }

	// // add groups before adding/updating users
	// err = handleGroupsResp(ctx, groupsResp)
	// if err != nil {
	// 	return err
	// }
	// err = handleUsersResp(ctx, usersResp)
	// if err != nil {
	// 	return err
	// }

	// err = handleLocked(ctx, lockedResp)
	// if err != nil {
	// 	return err
	// }

	log.Info("done", nil)
	return nil
}

// StartUpdater starts auth updater.
func (c Client) StartUpdater(ctx context.Context, updateCh <-chan struct{}) error {
	setupAuth(AuthLinux{})

	for {
		err := c.updateAll(ctx)
		if err != nil {
			return err
		}

		select {
		case <-updateCh:
			// received update notification
		case <-ctx.Done():
			return nil
		}
	}
}

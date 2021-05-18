package etcdpasswd

import (
	"context"
	"strconv"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/clientv3util"
)

// Client provides high-level API to edit etcd database.
type Client struct {
	*clientv3.Client
}

func (c Client) list(ctx context.Context, prefix string) ([]string, error) {
	resp, err := c.Get(ctx, prefix,
		clientv3.WithPrefix(),
		clientv3.WithKeysOnly(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return nil, err
	}

	ret := make([]string, resp.Count)
	for i, kv := range resp.Kvs {
		ret[i] = string(kv.Key[len(prefix):])
	}
	return ret, nil
}

func (c Client) getLastID(ctx context.Context, key string, startID int) (_ int, rev int64, e error) {
RETRY:
	resp, err := c.Get(ctx, key)
	if err != nil {
		e = err
		return
	}

	if resp.Count == 0 {
		v := strconv.Itoa(startID)
		_, err = c.Txn(ctx).
			If(clientv3util.KeyMissing(key)).
			Then(clientv3.OpPut(key, v)).
			Commit()
		if err != nil {
			e = err
			return
		}
		goto RETRY
	}

	uid, err := strconv.Atoi(string(resp.Kvs[0].Value))
	if err != nil {
		e = err
		return
	}

	return uid, resp.Kvs[0].ModRevision, nil
}

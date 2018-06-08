package etcdpasswd

import (
	"context"

	"github.com/coreos/etcd/clientv3"
)

// Client provides high-level API to edit etcd database.
type Client struct {
	*clientv3.Client
}

func (c Client) list(ctx context.Context, prefix string) ([]string, error) {
	resp, err := c.Get(ctx, prefix, clientv3.WithPrefix(),
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

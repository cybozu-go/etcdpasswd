package etcdpasswd

import (
	"github.com/coreos/etcd/clientv3"
)

// Client provides high-level API to edit etcd database.
type Client struct {
	*clientv3.Client
}

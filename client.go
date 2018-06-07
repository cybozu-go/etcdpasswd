package etcdpasswd

import (
	"errors"

	"github.com/coreos/etcd/clientv3"
)

// Client provides high-level API to edit etcd database.
type Client struct {
	*EtcdConfig
	*clientv3.Client
}

// ErrCASFailure indicates compare-and-swap failure
var ErrCASFailure = errors.New("CAS failed")

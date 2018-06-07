package cli

import "github.com/coreos/etcd/clientv3"

var (
	etcdClient *clientv3.Client
)

func Setup(client *clientv3.Client) {
	etcdClient = client
}

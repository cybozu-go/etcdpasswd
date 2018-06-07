package etcdpasswd

import (
	"path"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	defaultEtcdPrefix = "/passwd/"
)

// EtcdConfig represents configuration parameters to access etcd.
type EtcdConfig struct {
	Servers  []string `yaml:"servers"`
	Prefix   string   `yaml:"prefix"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

func NewEtcdConfig() *EtcdConfig {
	return &EtcdConfig{
		Prefix: defaultEtcdPrefix,
	}
}

func (ec *EtcdConfig) Key(p string) string {
	key := path.Join(ec.Prefix, p)
	if len(key) < ec.Prefix {
		return ec.Prefix
	}
	return key
}

func (ec *EtcdConfig) Client() (*clientv3.Client, error) {
	etcdCfg := clientv3.Config{
		Endpoints:   ec.Servers,
		DialTimeout: 2 * time.Second,
		Username:    ec.EtcdUsername,
		Password:    ec.EtcdPassword,
	}
	return clientv3.New(etcdCfg)
}

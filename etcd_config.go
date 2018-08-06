package etcdpasswd

import (
	"github.com/cybozu-go/etcdutil"
)

const (
	defaultEtcdPrefix = "/passwd/"
)

// NewEtcdConfig creates Config with default prefix.
func NewEtcdConfig() *etcdutil.Config {
	return etcdutil.NewConfig(defaultEtcdPrefix)
}

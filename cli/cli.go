package cli

import "github.com/cybozu-go/etcdpasswd"

var client etcdpasswd.Client

// Setup setups this package.
func Setup(c etcdpasswd.Client) {
	client = c
}

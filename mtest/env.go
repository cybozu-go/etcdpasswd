package mtest

import (
	"os"
)

var (
	bridgeAddress = os.Getenv("BRIDGE_ADDRESS")
	host1         = os.Getenv("HOST1")
	host2         = os.Getenv("HOST2")
	host3         = os.Getenv("HOST3")

	etcdPath       = os.Getenv("ETCD")
	etcdctlPath    = os.Getenv("ETCDCTL")
	etcdpasswdPath = os.Getenv("ETCDPASSWD")
	epagentPath    = os.Getenv("EPAGENT")

	sshKeyFile = os.Getenv("SSH_PRIVKEY")
)

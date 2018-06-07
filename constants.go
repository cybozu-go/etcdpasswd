package etcdpasswd

// Internal schema keys.
const (
	KeyConfig  = "config"
	KeyLastUID = "last-uid"
	KeyLastGID = "last-gid"
	KeyUsers   = "users"
	KeyGroups  = "groups"
	KeyLocked  = "locked"
)

const (
	defaultShell = "/bin/bash"
)

package etcdpasswd

// Internal schema keys.
const (
	KeyConfig        = "config"
	KeyLastUID       = "last-uid"
	KeyLastGID       = "last-gid"
	KeyUsers         = "users"
	KeyDeletedUsers  = "deleted-users"
	KeyGroups        = "groups"
	KeyDeletedGroups = "deleted-groups"
	KeyLocked        = "locked"
)

const (
	// DefaultShell is the default shell program.
	DefaultShell = "/bin/bash"
)

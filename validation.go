package etcdpasswd

import "regexp"

var (
	reValidUser  = regexp.MustCompile(`^[a-z][-a-z0-9_]*$`)
	reValidGroup = regexp.MustCompile(`^[a-z][-a-z0-9_]*$`)
)

// IsValidUserName returns true if name is valid for etcdpasswd managed user.
func IsValidUserName(name string) bool {
	switch name {
	case "root", "nobody":
		return false
	}
	return reValidUser.MatchString(name)
}

// IsValidGroupName returns true if name is valid for etcdpasswd managed group.
func IsValidGroupName(name string) bool {
	switch name {
	case "root", "nogroup", "adm", "sudo":
		return false
	}
	return reValidGroup.MatchString(name)
}

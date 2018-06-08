package etcdpasswd

//go:generate mockgen -destination=mocks/auth.go -package=mocks github.com/cybozu-go/etcdpasswd Auth

import (
	"context"
	"os/user"

	"github.com/cybozu-go/cmd"
)

// Auth is an interface for the target of auth operations.
type Auth interface {
	DeleteUser(context.Context, string) error
}

var auth Auth

func setupAuth(a Auth) {
	auth = a
}

// AuthLinux implements Auth for a real Linux server.
type AuthLinux struct {
}

// DeleteUser deletes a Linux user.
func (a AuthLinux) DeleteUser(ctx context.Context, username string) error {
	_, err := user.Lookup(username)
	if err != nil {
		if _, ok := err.(user.UnknownUserError); ok {
			// already deleted
			return nil
		}
		return err
	}

	return cmd.CommandContext(ctx, "userdel", "--force", username).Run()
}

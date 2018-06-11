package cli

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/google/subcommands"
	"golang.org/x/crypto/ssh"
)

type cert struct{}

func (c cert) SetFlags(f *flag.FlagSet) {}

func (c cert) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	newc := NewCommander(f, "cert")
	newc.Register(certListCommand(), "")
	newc.Register(certAddCommand(), "")
	newc.Register(certRemoveCommand(), "")
	return newc.Execute(ctx)
}

// CertCommand implements "cert" subcommand.
func CertCommand() subcommands.Command {
	return subcmd{
		cert{},
		"cert",
		"manage SSH public keys of users",
		"cert ACTION ...",
	}
}

func pprintPubKey(pubkey string) string {
	pk, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(pubkey))
	if err != nil {
		return pubkey
	}
	return fmt.Sprintf("%s (%s)", comment, pk.Type())
}

type certList struct{}

func (c certList) SetFlags(f *flag.FlagSet) {}

func (c certList) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 1 {
		f.Usage()
		return subcommands.ExitUsageError
	}

	name := f.Arg(0)
	user, _, err := client.GetUser(ctx, name)
	if err != nil {
		return handleError(err)
	}

	for i, pubkey := range user.PubKeys {
		fmt.Printf("%d: %s\n", i, pprintPubKey(pubkey))
	}
	return handleError(nil)
}

func certListCommand() subcommands.Command {
	return subcmd{
		certList{},
		"list",
		"list SSH public keys of a user",
		"cert list NAME",
	}
}

type certAdd struct{}

func (c certAdd) SetFlags(f *flag.FlagSet) {}

func (c certAdd) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	var name string
	var file string
	switch f.NArg() {
	case 2:
		file = f.Arg(1)
		fallthrough
	case 1:
		name = f.Arg(0)
	default:
		f.Usage()
		return subcommands.ExitUsageError
	}

	input := os.Stdin
	if len(file) > 0 {
		g, err := os.Open(file)
		if err != nil {
			return handleError(err)
		}
		defer g.Close()
		input = g
	}
	pubkey, err := ioutil.ReadAll(input)
	if err != nil {
		return handleError(err)
	}

	user, rev, err := client.GetUser(ctx, name)
	if err != nil {
		return handleError(err)
	}

	user.PubKeys = append(user.PubKeys, string(bytes.TrimSpace(pubkey)))
	err = client.UpdateUser(ctx, user, rev)
	return handleError(err)
}

func certAddCommand() subcommands.Command {
	return subcmd{
		certAdd{},
		"add",
		"add a SSH public key of a user",
		`cert add NAME [FILE]

If FILE is not specified, public key is read from stdin.`,
	}
}

type certRemove struct{}

func (c certRemove) SetFlags(f *flag.FlagSet) {}

func (c certRemove) Execute(ctx context.Context, f *flag.FlagSet) subcommands.ExitStatus {
	if f.NArg() != 2 {
		f.Usage()
		return subcommands.ExitUsageError
	}

	name := f.Arg(0)
	index, err := strconv.Atoi(f.Arg(1))
	if err != nil {
		return handleError(err)
	}

	user, rev, err := client.GetUser(ctx, name)
	if err != nil {
		return handleError(err)
	}

	if index >= len(user.PubKeys) || index < 0 {
		return handleError(errors.New("invalid index"))
	}

	user.PubKeys = append(user.PubKeys[:index], user.PubKeys[(index+1):]...)
	err = client.UpdateUser(ctx, user, rev)
	return handleError(err)
}

func certRemoveCommand() subcommands.Command {
	return subcmd{
		certRemove{},
		"remove",
		"remove a SSH public key of a user",
		"cert remove NAME INDEX",
	}
}

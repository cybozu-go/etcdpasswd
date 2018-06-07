package main

import (
	"context"
	"flag"
	"os"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/etcdpasswd/cli"
	"github.com/cybozu-go/log"
	"github.com/google/subcommands"
	yaml "gopkg.in/yaml.v2"
)

var (
	flgConfigPath = flag.String("config", "/etc/etcdpasswd.yml", "configuration file path")
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "misc")
	subcommands.Register(subcommands.FlagsCommand(), "misc")
	subcommands.Register(subcommands.CommandsCommand(), "misc")
	subcommands.Register(cli.SetCommand(), "")
	subcommands.Register(cli.GetCommand(), "")
	subcommands.Register(cli.LockCommand(), "")
	subcommands.Register(cli.UnlockCommand(), "")
	subcommands.Register(cli.UserCommand(), "")
	subcommands.Register(cli.CertCommand(), "")
	subcommands.Register(cli.GroupCommand(), "")
	flag.Parse()
	cmd.LogConfig{}.Apply()

	f, err := os.Open(*flgConfigPath)
	if err != nil {
		log.ErrorExit(err)
	}
	defer f.Close()

	etcdConfig := etcdpasswd.NewEtcdConfig()
	err = yaml.NewDecoder(f).Decode(etcdConfig)
	if err != nil {
		log.ErrorExit(err)
	}

	client, err := etcdConfig.Client()
	if err != nil {
		log.ErrorExit(err)
	}
	defer client.Close()

	cli.Setup(client)

	exitStatus := subcommands.ExitSuccess
	cmd.Go(func(ctx context.Context) error {
		exitStatus = subcommands.Execute(ctx)
		return nil
	})
	cmd.Stop()
	cmd.Wait()
	os.Exit(int(exitStatus))
}

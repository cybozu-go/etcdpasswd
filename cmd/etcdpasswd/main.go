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

func loadConfig(p string) (*etcdpasswd.EtcdConfig, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := etcdpasswd.NewEtcdConfig()
	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "misc")
	subcommands.Register(subcommands.FlagsCommand(), "misc")
	subcommands.Register(subcommands.CommandsCommand(), "misc")
	subcommands.Register(cli.SetCommand(), "")
	subcommands.Register(cli.GetCommand(), "")
	subcommands.Register(cli.LockCommand(), "")
	subcommands.Register(cli.UnlockCommand(), "")
	subcommands.Register(cli.UserCommand(), "")
	// subcommands.Register(cli.CertCommand(), "")
	// subcommands.Register(cli.GroupCommand(), "")
	flag.Parse()
	cmd.LogConfig{}.Apply()

	cfg, err := loadConfig(*flgConfigPath)
	if err != nil {
		flag.Usage()
		log.ErrorExit(err)
	}

	etcd, err := cfg.Client()
	if err != nil {
		log.ErrorExit(err)
	}
	defer etcd.Close()

	client := etcdpasswd.Client{
		Client: etcd,
	}
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

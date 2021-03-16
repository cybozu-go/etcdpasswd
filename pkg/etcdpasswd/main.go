package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/etcdpasswd/cli"
	"github.com/cybozu-go/etcdutil"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	"github.com/google/subcommands"
	"sigs.k8s.io/yaml"
)

var (
	flgConfigPath = flag.String("config", "/etc/etcdpasswd/config.yml", "configuration file path")
	flgVersion    = flag.Bool("version", false, "version")
)

func loadConfig(p string) (*etcdutil.Config, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	cfg := etcdpasswd.NewEtcdConfig()
	err = yaml.Unmarshal(b, cfg)
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
	subcommands.Register(cli.UserCommand(), "")
	subcommands.Register(cli.CertCommand(), "")
	subcommands.Register(cli.GroupCommand(), "")
	subcommands.Register(cli.LockerCommand(), "")
	flag.Parse()
	well.LogConfig{}.Apply()

	if *flgVersion {
		fmt.Println(etcdpasswd.Version)
		os.Exit(0)
	}

	cfg, err := loadConfig(*flgConfigPath)
	if err != nil {
		flag.Usage()
		log.ErrorExit(err)
	}

	etcd, err := etcdutil.NewClient(cfg)
	if err != nil {
		log.ErrorExit(err)
	}
	defer etcd.Close()

	client := etcdpasswd.Client{
		Client: etcd,
	}
	cli.Setup(client)

	exitStatus := subcommands.ExitSuccess
	well.Go(func(ctx context.Context) error {
		exitStatus = subcommands.Execute(ctx)
		return nil
	})
	well.Stop()
	well.Wait()
	os.Exit(int(exitStatus))
}

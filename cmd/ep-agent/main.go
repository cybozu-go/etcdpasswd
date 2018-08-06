package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/etcdpasswd/agent"
	"github.com/cybozu-go/etcdpasswd/syncer"
	"github.com/cybozu-go/etcdutil"
	"github.com/cybozu-go/log"
	yaml "gopkg.in/yaml.v2"
)

var (
	flgConfigPath = flag.String("config", "/etc/etcdpasswd.yml", "configuration file path")
	flgSyncer     = flag.String("syncer", "os", "user sync driver [os,dummy]")
)

func loadConfig(p string) (*etcdutil.Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := etcdutil.NewConfig(etcdpasswd.DefaultEtcdPrefix)
	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	flag.Parse()
	cmd.LogConfig{}.Apply()

	cfg, err := loadConfig(*flgConfigPath)
	if err != nil {
		log.ErrorExit(err)
	}

	etcd, err := etcdutil.NewClient(cfg)
	if err != nil {
		log.ErrorExit(err)
	}
	defer etcd.Close()

	var sc etcdpasswd.Syncer
	switch *flgSyncer {
	case "os":
		sc = syncer.UbuntuSyncer{}
	case "dummy":
		sc = syncer.NewDummySyncer()
	default:
		fmt.Fprintln(os.Stderr, "no such syncer: "+*flgSyncer)
		os.Exit(1)
	}

	agent := agent.Agent{Client: etcd, Syncer: sc}

	updateCh := make(chan struct{}, 1)
	cmd.Go(func(ctx context.Context) error {
		return agent.StartWatching(ctx, updateCh)
	})
	cmd.Go(func(ctx context.Context) error {
		return agent.StartUpdater(ctx, updateCh)
	})
	cmd.Stop()

	err = cmd.Wait()
	if err != nil && !cmd.IsSignaled(err) {
		log.ErrorExit(err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/etcdpasswd/agent"
	"github.com/cybozu-go/etcdpasswd/syncer"
	"github.com/cybozu-go/etcdutil"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	"sigs.k8s.io/yaml"
)

var (
	flgConfigPath = flag.String("config", "/etc/etcdpasswd/config.yml", "configuration file path")
	flgSyncer     = flag.String("syncer", "os", "user sync driver [os,dummy]")
	flgVersion    = flag.Bool("version", false, "version")
)

func loadConfig(p string) (*etcdutil.Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
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
	flag.Parse()
	well.LogConfig{}.Apply()

	if *flgVersion {
		fmt.Println(etcdpasswd.Version)
		os.Exit(0)
	}

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
	well.Go(func(ctx context.Context) error {
		return agent.StartWatching(ctx, updateCh)
	})
	well.Go(func(ctx context.Context) error {
		return agent.StartUpdater(ctx, updateCh)
	})
	well.Stop()

	err = well.Wait()
	if err != nil && !well.IsSignaled(err) {
		log.ErrorExit(err)
	}
}

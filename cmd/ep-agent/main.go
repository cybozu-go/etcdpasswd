package main

import (
	"context"
	"flag"
	"os"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/etcdpasswd/cli"
	"github.com/cybozu-go/log"
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
	flag.Parse()
	cmd.LogConfig{}.Apply()

	cfg, err := loadConfig(*flgConfigPath)
	if err != nil {
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

	updateCh := make(chan struct{}, 1)
	cmd.Go(func(ctx context.Context) error {
		return client.StartWatching(ctx, updateCh)
	})
	cmd.Go(func(ctx context.Context) error {
		return client.StartUpdater(ctx, updateCh)
	})
	cmd.Stop()

	err = cmd.Wait()
	if !cmd.IsSignaled(err) && err != nil {
		log.ErrorExit(err)
	}
}

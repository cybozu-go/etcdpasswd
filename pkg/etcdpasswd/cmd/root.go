package cmd

import (
	"os"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/etcdutil"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	"github.com/spf13/cobra"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sigs.k8s.io/yaml"
)

var (
	cfgFile    string
	etcdClient *clientv3.Client
	client     etcdpasswd.Client
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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "etcdpasswd",
	Short:   "CLI tool to edit the central database on etcd",
	Long:    "CLI tool to edit the central database on etcd.",
	Version: etcdpasswd.Version,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// without this, each subcommand's RunE would display usage text.
		cmd.SilenceUsage = true

		err := well.LogConfig{}.Apply()
		if err != nil {
			return err
		}

		cfg, err := loadConfig(cfgFile)
		if err != nil {
			return err
		}

		etcd, err := etcdutil.NewClient(cfg)
		if err != nil {
			return err
		}
		etcdClient = etcd

		client = etcdpasswd.Client{
			Client: etcd,
		}

		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if etcdClient != nil {
			etcdClient.Close()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.ErrorExit(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/etcdpasswd/config.yml", "configuration file path")
}

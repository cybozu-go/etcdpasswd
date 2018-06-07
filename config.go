package etcdpasswd

import (
	"context"
	"encoding/json"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/clientv3util"
)

// Config represents etcdpasswd configurations
type Config struct {
	StartUID      int      `json:"start-uid"`
	StartGID      int      `json:"start-gid"`
	DefaultGroup  string   `json:"default-group"`
	DefaultGroups []string `json:"default-groups"`
	DefaultShell  string   `json:"default-shell"`
}

func (c Client) initializeConfig(ctx context.Context, key string) error {
	config := &Config{
		DefaultShell: defaultShell,
	}
	j, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = c.Txn(ctx).
		If(clientv3util.KeyMissing(key)).
		Then(clientv3.OpPut(key, string(j))).
		Commit()

	return err
}

// GetConfig retrieves *Config with revision.
func (c Client) GetConfig(ctx context.Context) (*Config, int64, error) {
	key := c.Key(KeyConfig)

RETRY:
	resp, err := c.Get(ctx, key)
	if err != nil {
		return nil, 0, err
	}

	if resp.Count == 0 {
		err = c.initializeConfig(ctx, key)
		if err != nil {
			return nil, 0, err
		}
		goto RETRY
	}

	config := new(Config)
	err = json.Unmarshal(resp.Kvs[0].Value, config)
	if err != nil {
		return nil, 0, err
	}
	return config, resp.Kvs[0].ModRevision, nil
}

// SetConfig tries to update *Config.
// If update was conflicted, ErrCASFailure is returned.
func (c Client) SetConfig(ctx context.Context, cfg *Config, rev int64) error {
	key := c.Key(KeyConfig)

	j, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	resp, err := c.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(key), "=", rev)).
		Then(clientv3.OpPut(key, string(j))).
		Commit()
	if err != nil {
		return err
	}
	if !resp.Succeeded {
		return ErrCASFailure
	}

	return nil
}

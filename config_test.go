package etcdpasswd

import (
	"context"
	"reflect"
	"testing"
)

func testConfigGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	config, rev, err := client.GetConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if config == nil {
		t.Fatal("config must not be nil")
	}

	if config.DefaultShell != defaultShell {
		t.Error("wrong default shell:", config.DefaultShell)
	}
	if config.StartUID != 0 {
		t.Error("StartUID must be 0")
	}
	if config.StartGID != 0 {
		t.Error("StartGID must be 0")
	}

	j := `
{
  "start-uid": 100,
  "default-groups": ["cybozu"]
}
`
	_, err = client.Put(ctx, client.Key(KeyConfig), j)
	if err != nil {
		t.Fatal(err)
	}

	config, rev2, err := client.GetConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if config == nil {
		t.Fatal("config must not be nil")
	}

	if rev == rev2 {
		t.Error("revision not updated")
	}
	if config.StartUID != 100 {
		t.Error("wrong StartUID:", config.StartUID)
	}
	if !reflect.DeepEqual(config.DefaultGroups, []string{"cybozu"}) {
		t.Error("wrong DefaultGroups", config.DefaultGroups)
	}
}

func testConfigSet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	config, rev, err := client.GetConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if config == nil {
		t.Fatal("config must not be nil")
	}

	config.StartUID = 123
	err = client.SetConfig(ctx, config, rev)
	if err != nil {
		t.Fatal(err)
	}

	config, rev, err = client.GetConfig(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if config == nil {
		t.Fatal("config must not be nil")
	}

	if config.StartUID != 123 {
		t.Error("StartUID not updated:", config.StartUID)
	}

	// increment revision by hand
	j := `
{
  "start-uid": 100,
  "default-groups": ["cybozu"]
}
`
	_, err = client.Put(ctx, client.Key(KeyConfig), j)
	if err != nil {
		t.Fatal(err)
	}

	// try to set with old rev
	err = client.SetConfig(ctx, config, rev)
	if err != ErrCASFailure {
		t.Error("CAS should have failed:", err)
	}
}

func TestConfig(t *testing.T) {
	t.Run("Get", testConfigGet)
	t.Run("Set", testConfigSet)
}

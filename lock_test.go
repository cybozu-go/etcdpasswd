package etcdpasswd

import (
	"context"
	"path"
	"testing"

	"github.com/coreos/etcd/clientv3"
)

func TestLock(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	err := client.Lock(ctx, "root")
	if err == nil {
		t.Error("root should not be locked")
	}

	err = client.Lock(ctx, "cybozu")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Get(ctx, path.Join(KeyLocked, "cybozu"), clientv3.WithCountOnly())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 1 {
		t.Error("cybozu was not added to locked user database")
	}

	err = client.Lock(ctx, "cybozu")
	if err != nil {
		t.Error("failed to daouble lock cybozu")
	}
}

func TestUnlock(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	err := client.Unlock(ctx, "cybozu")
	if err != nil {
		t.Error("failed to unlock non-locked user")
	}

	_, err = client.Put(ctx, path.Join(KeyLocked, "cybozu"), "")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Unlock(ctx, "cybozu")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Get(ctx, path.Join(KeyLocked, "cybozu"), clientv3.WithCountOnly())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 0 {
		t.Error("unlock did not remove user from locked user database")
	}
}

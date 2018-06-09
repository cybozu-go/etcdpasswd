package etcdpasswd

import (
	"context"
	"path"
	"reflect"
	"testing"

	"github.com/coreos/etcd/clientv3"
)

func testGroupPrepare(t *testing.T, client Client) {
	ctx := context.Background()
	_, err := client.Put(ctx, path.Join(KeyGroups, "test1"), "12345")
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Put(ctx, path.Join(KeyGroups, "test2"), "12346")
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Put(ctx, path.Join(KeyGroups, "abc"), "12347")
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Put(ctx, path.Join(KeyDeletedGroups, "test0"), "12344")
	if err != nil {
		t.Fatal(err)
	}
}

func testGroupList(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	testGroupPrepare(t, client)

	groups, err := client.ListGroups(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 3 {
		t.Fatal(`len(groups) != 3`, len(groups))
	}

	expected := []Group{
		{"abc", 12347}, {"test1", 12345}, {"test2", 12346},
	}
	if !reflect.DeepEqual(expected, groups) {
		t.Error(`!reflect.DeepEqual(expected, groups)`, groups)
	}
}

func testGroupAdd(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	testGroupPrepare(t, client)

	err := client.AddGroup(ctx, "test0")
	if err == nil {
		t.Fatal("group should not be added w/o proper start-gid")
	}

	_, err = client.Put(ctx, KeyConfig, testConfigData)
	if err != nil {
		t.Fatal(err)
	}

	err = client.AddGroup(ctx, "test1")
	if err != ErrExists {
		t.Error(`err != ErrExists`)
	}

	err = client.AddGroup(ctx, "test0")
	if err != nil {
		t.Fatal(err)
	}
	delKey := path.Join(KeyDeletedGroups, "test0")
	resp, err := client.Get(ctx, delKey, clientv3.WithCountOnly())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 0 {
		t.Error(`resp.Count != 0`, resp.Count)
	}
	resp2, err := client.Get(ctx, path.Join(KeyGroups, "test0"))
	if err != nil {
		t.Fatal(err)
	}
	if resp2.Count != 1 {
		t.Fatal(`resp2.Count != 1`)
	}
	if string(resp2.Kvs[0].Value) != "3000" {
		t.Error(string(resp2.Kvs[0].Value) != "3000", string(resp2.Kvs[0].Value))
	}

	err = client.AddGroup(ctx, "zzz")
	if err != nil {
		t.Fatal(err)
	}
	resp3, err := client.Get(ctx, path.Join(KeyGroups, "zzz"))
	if err != nil {
		t.Fatal(err)
	}
	if resp3.Count != 1 {
		t.Fatal(`resp3.Count != 1`)
	}
	if string(resp3.Kvs[0].Value) != "3001" {
		t.Error(string(resp3.Kvs[0].Value) != "3001", string(resp3.Kvs[0].Value))
	}
}

func testGroupRemove(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	testGroupPrepare(t, client)

	err := client.RemoveGroup(ctx, "cybozu")
	if err != ErrNotFound {
		t.Fatal("non-existing group should not be removed")
	}

	err = client.RemoveGroup(ctx, "test1")
	if err != nil {
		t.Error(err)
	}

	delKey := path.Join(KeyDeletedGroups, "test1")
	resp, err := client.Get(ctx, delKey, clientv3.WithCountOnly())
	if resp.Count != 1 {
		t.Error(`deleted group should have been registered`)
	}
}

func TestGroup(t *testing.T) {
	t.Run("List", testGroupList)
	t.Run("Add", testGroupAdd)
	t.Run("Remove", testGroupRemove)
}

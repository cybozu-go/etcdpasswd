package etcdpasswd

import (
	"context"
	"reflect"
	"testing"

	"github.com/coreos/etcd/clientv3"
)

var testUserData = `
{
    "name": "cybozu",
    "uid": 2001,
    "display-name": "Cy Bozu",
    "group": "cybozu",
    "groups": ["sudo", "adm"],
    "shell": "/bin/zsh",
    "public-keys": ["aaa", "bbb"]
}
`

var testConfigData = `
{
    "start-uid": 2000,
    "start-gid": 3000,
    "default-shell": "/bin/bash"
}
`

var testConfigData2 = `
{
    "start-uid": 2000,
    "start-gid": 3000,
    "default-group": "cybozu",
    "default-groups": ["sudo", "adm"],
    "default-shell": "/bin/bash"
}
`

func testUserGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	_, _, err := client.GetUser(ctx, "cybozu")
	if err != ErrNotFound {
		t.Error("user cybozu should not be found")
	}

	key := KeyUsers + "cybozu"
	_, err = client.Put(ctx, key, testUserData)
	if err != nil {
		t.Fatal(err)
	}

	u, _, err := client.GetUser(ctx, "cybozu")
	if err != nil {
		t.Fatal(err)
	}
	if u.Name != "cybozu" {
		t.Error("name must be cybozu,", u.Name)
	}
}

func testUserAdd(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	u := &User{
		Name:        "cybozu",
		DisplayName: "Cy Bozu",
		Shell:       "/bin/zsh",
	}

	err := client.AddUser(ctx, u)
	if err == nil {
		t.Fatal("user should not be added w/o proper start-uid")
	}

	_, err = client.Put(ctx, KeyConfig, testConfigData)
	if err != nil {
		t.Fatal(err)
	}

	err = client.AddUser(ctx, u)
	if err == nil {
		t.Fatal("user should not be added w/o group")
	}

	u2 := *u
	u2.Name = "cybozu2"
	u2.Group = "foo"
	err = client.AddUser(ctx, &u2)
	if err != nil {
		t.Fatal(err)
	}

	gu2, _, err := client.GetUser(ctx, "cybozu2")
	if err != nil {
		t.Fatal(err)
	}
	if gu2.UID != 2000 {
		t.Error(`gu2.UID != 2000`, gu2.UID)
	}
	if gu2.Group != "foo" {
		t.Error(`gu2.Group != "foo"`)
	}

	_, err = client.Put(ctx, KeyConfig, testConfigData2)
	if err != nil {
		t.Fatal(err)
	}

	err = client.AddUser(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	err = client.AddUser(ctx, u)
	if err != ErrExists {
		t.Error(`err != ErrExists`)
	}

	gu, _, err := client.GetUser(ctx, "cybozu")
	if err != nil {
		t.Fatal(err)
	}

	if gu.Name != "cybozu" {
		t.Error(`gu.Name != "cybozu"`)
	}
	if gu.UID != 2001 {
		t.Error(`gu.UID != 2001`)
	}
	if gu.Group != "cybozu" {
		t.Error(`gu.Name != "cybozu"`)
	}
	if !reflect.DeepEqual(gu.Groups, []string{"sudo", "adm"}) {
		t.Error(`!reflect.DeepEqual(gu.Groups, []string{"sudo", "adm"})`)
	}
	if gu.Shell != "/bin/zsh" {
		t.Error(`gu.Shell != "/bin/zsh"`)
	}

	delKey := KeyDeletedUsers + "cybozu3"
	_, err = client.Put(ctx, delKey, "")
	if err != nil {
		t.Fatal(err)
	}

	u3 := *u
	u3.Name = "cybozu3"
	u3.Groups = []string{"foo"}
	u3.Shell = ""
	err = client.AddUser(ctx, &u3)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Get(ctx, delKey, clientv3.WithCountOnly())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 0 {
		t.Error("entry for deleted-users should have been removed")
	}

	gu3, _, err := client.GetUser(ctx, "cybozu3")
	if err != nil {
		t.Fatal(err)
	}
	if gu3.UID != 2002 {
		t.Error(`gu3.UID != 2002`)
	}
	if !reflect.DeepEqual(gu3.Groups, []string{"foo"}) {
		t.Error(`!reflect.DeepEqual(gu3.Groups, []string{"foo"})`)
	}
	if gu3.Shell != "/bin/bash" {
		t.Error(`gu3.Shell != "/bin/bash", `, gu3.Shell)
	}
}

func testUserUpdate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	key := KeyUsers + "cybozu"
	_, err := client.Put(ctx, key, testUserData)
	if err != nil {
		t.Fatal(err)
	}

	u, rev, err := client.GetUser(ctx, "cybozu")
	if err != nil {
		t.Fatal(err)
	}

	nu := *u
	nu.Name = "notexist"
	err = client.UpdateUser(ctx, &nu, rev)
	if err != ErrCASFailure {
		t.Fatal("non-existing user should not be updated")
	}

	nu = *u
	nu.Group = ""
	err = client.UpdateUser(ctx, &nu, rev)
	if err == nil {
		t.Fatal("empty group should not be allowed")
	}

	nu = *u
	nu.Shell = ""
	err = client.UpdateUser(ctx, &nu, rev)
	if err == nil {
		t.Fatal("empty shell should not be allowed")
	}

	nu = *u
	nu.Group = "cybozu2"
	err = client.UpdateUser(ctx, &nu, rev)
	if err != nil {
		t.Error(err)
	}

	gu, _, err := client.GetUser(ctx, "cybozu")
	if err != nil {
		t.Fatal(err)
	}
	if gu.Group != "cybozu2" {
		t.Error(`gu.Group != "cybozu2"`)
	}

	nu = *u
	nu.Group = "cybozu3"
	err = client.UpdateUser(ctx, &nu, rev)
	if err != ErrCASFailure {
		t.Error(`err != ErrCASFailure`)
	}
}

func testUserRemove(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	err := client.RemoveUser(ctx, "cybozu")
	if err != ErrNotFound {
		t.Fatal("non-existing user should not be removed")
	}

	key := KeyUsers + "cybozu"
	_, err = client.Put(ctx, key, testUserData)
	if err != nil {
		t.Fatal(err)
	}

	err = client.RemoveUser(ctx, "cybozu")
	if err != nil {
		t.Error(err)
	}

	delKey := KeyDeletedUsers + "cybozu"
	resp, _ := client.Get(ctx, delKey, clientv3.WithCountOnly())
	if resp.Count != 1 {
		t.Error(`deleted user should have been registered`)
	}
}

func TestUser(t *testing.T) {
	t.Run("Get", testUserGet)
	t.Run("Add", testUserAdd)
	t.Run("Update", testUserUpdate)
	t.Run("Remove", testUserRemove)
}

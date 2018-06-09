package etcdpasswd

import (
	"context"
	"path"
	"testing"
)

func testPrepareDatabase(ctx context.Context, t *testing.T, client Client) int64 {
	_, err := client.Put(ctx, KeyConfig, testConfigData)
	if err != nil {
		t.Fatal(err)
	}

	err = client.AddUser(ctx, &User{
		Name:        "user1",
		DisplayName: "u ser 1",
		Group:       "group1",
		Groups:      []string{"group2", "group3"},
		Shell:       "/bin/sh",
		PubKeys:     []string{"pubkey1"},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.AddUser(ctx, &User{
		Name:        "user2",
		DisplayName: "u ser 2",
		Group:       "group1",
		Shell:       "/bin/sh",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.AddGroup(ctx, "group1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Put(ctx, path.Join(KeyDeletedUsers, "user3"), "")
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.Put(ctx, path.Join(KeyDeletedGroups, "group2"), "")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Put(ctx, path.Join(KeyLocked, "user4"), "")
	if err != nil {
		t.Fatal(err)
	}

	return resp.Header.Revision
}

func testGetDatabase(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	rev := testPrepareDatabase(ctx, t, client)

	err := client.AddUser(ctx, &User{
		Name:  "user3",
		Group: "group3",
	})
	if err != nil {
		t.Fatal(err)
	}

	// get the latest snapshot
	db, err := GetDatabase(ctx, client.Client, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(db.Users) != 3 {
		t.Fatal(`len(db.Users) != 3`, len(db.Users))
	}
	if db.Users[0].Name != "user1" {
		t.Error(`db.Users[0].Name != "user1"`)
	}
	if db.Users[1].Name != "user2" {
		t.Error(`db.Users[1].Name != "user2"`)
	}
	if db.Users[2].Name != "user3" {
		t.Error(`db.Users[2].Name != "user3"`)
	}
	if len(db.Groups) != 1 {
		t.Fatal(`len(db.Groups) != 1`, len(db.Groups))
	}
	if db.Groups[0].Name != "group1" {
		t.Error(`db.Groups[0].Name != "group1"`)
	}
	if len(db.DeletedUsers) != 0 {
		t.Error(`len(db.DeletedUsers) != 0`)
	}
	if len(db.DeletedGroups) != 1 {
		t.Fatal(`len(db.DeletedGroups) != 1`)
	}
	if db.DeletedGroups[0] != "group2" {
		t.Error(`db.DeletedGroups[0] != "group2"`)
	}
	if len(db.LockedUsers) != 1 {
		t.Fatal(`len(db.LockedUsers) != 1`)
	}
	if db.LockedUsers[0] != "user4" {
		t.Error(`db.LockedUsers[0] != "user4"`)
	}

	// get a previous snapshot
	db, err = GetDatabase(ctx, client.Client, rev)
	if err != nil {
		t.Fatal(err)
	}
	if len(db.Users) != 2 {
		t.Fatal(`len(db.Users) != 2`, len(db.Users))
	}
	if len(db.DeletedUsers) != 1 {
		t.Fatal(`len(db.DeletedUsers) != 1`)
	}
	if db.DeletedUsers[0] != "user3" {
		t.Error(`db.DeletedUsers[0] != "user3"`)
	}
	if len(db.LockedUsers) != 1 {
		t.Fatal(`len(db.LockedUsers) != 1`)
	}
	if db.LockedUsers[0] != "user4" {
		t.Error(`db.LockedUsers[0] != "user4"`)
	}
}

func TestDatabase(t *testing.T) {
	t.Run("Get", testGetDatabase)
}

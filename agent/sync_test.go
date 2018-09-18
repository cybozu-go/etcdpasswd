package agent

import (
	"context"
	"reflect"
	"testing"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/etcdpasswd/syncer"
)

func TestEqualStringSlice(t *testing.T) {
	t.Parallel()

	if equalStringSlice([]string{"abc"}, nil) {
		t.Error(`equalStringSlice([]string{"abc"}, nil)`)
	}
	if equalStringSlice([]string{"abc", "def"}, []string{"abc", "xyz"}) {
		t.Error(`equalStringSlice([]string{"abc", "def"}, []string{"abc", "xyz"})`)
	}
	if !equalStringSlice([]string{"abc", "def"}, []string{"abc", "def"}) {
		t.Error(`!equalStringSlice([]string{"abc", "def"}, []string{"abc", "def"})`)
	}
	if !equalStringSlice([]string{"def", "abc", "xyz"}, []string{"xyz", "def", "abc"}) {
		t.Error(`!equalStringSlice([]string{"def", "abc", "xyz"}, []string{"xyz", "def", "abc"})`)
	}
}

type User = etcdpasswd.User
type Group = etcdpasswd.Group

func TestSynchronize(t *testing.T) {
	t.Parallel()

	db := &etcdpasswd.Database{
		Users: []*User{
			{
				Name:        "user1",
				UID:         2000,
				DisplayName: "u ser 1",
				Group:       "system-group1",
				Groups:      []string{"group1", "group3", "system-group2"},
				Shell:       "/bin/sh",
				PubKeys:     []string{"pubkey1"},
			},
			{
				Name:        "user2",
				UID:         2001,
				DisplayName: "u ser 2",
				Group:       "group1",
				Shell:       "/bin/sh",
			},
		},
		Groups: []Group{
			{Name: "group1", GID: 3000},
			{Name: "group3", GID: 3002},
		},
		DeletedUsers:  []string{"user3"},
		DeletedGroups: []string{"group2"},
		LockedUsers:   []string{"system-user1", "system-user2"},
	}

	sc := syncer.NewMockSyncer()
	ctx := context.Background()

	err := sc.AddGroup(ctx, Group{Name: "system-group1", GID: 100})
	if err != nil {
		t.Fatal(err)
	}
	err = sc.AddGroup(ctx, Group{Name: "system-group2", GID: 101})
	if err != nil {
		t.Fatal(err)
	}
	err = sc.AddUser(ctx, &User{
		Name:  "system-user1",
		UID:   100,
		Group: "system-group1",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := syncer.NewMockSyncer()
	expected.Users["system-user1"] = &User{
		Name:  "system-user1",
		UID:   100,
		Group: "system-group1",
	}
	expected.Users["user1"] = &User{
		Name:        "user1",
		UID:         2000,
		DisplayName: "u ser 1",
		Group:       "system-group1",
		Groups:      []string{"group1", "group3", "system-group2"},
		Shell:       "/bin/sh",
		PubKeys:     []string{"pubkey1"},
	}
	expected.Users["user2"] = &User{
		Name:        "user2",
		UID:         2001,
		DisplayName: "u ser 2",
		Group:       "group1",
		Shell:       "/bin/sh",
	}
	expected.Groups["system-group1"] = &Group{Name: "system-group1", GID: 100}
	expected.Groups["system-group2"] = &Group{Name: "system-group2", GID: 101}
	expected.Groups["group1"] = &Group{Name: "group1", GID: 3000}
	expected.Groups["group3"] = &Group{Name: "group3", GID: 3002}
	expected.LockedUsers["system-user1"] = true

	// initial sync
	err = synchronize(ctx, db, sc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc, expected) {
		t.Errorf(`!reflect.DeepEqual(sc, expected): %#v`, sc)
	}

	// add a user
	db.Users = append(db.Users, &User{
		Name:  "user3",
		UID:   2004,
		Group: "group3",
	})
	db.DeletedUsers = db.DeletedUsers[:0]
	expected.Users["user3"] = &User{
		Name:  "user3",
		UID:   2004,
		Group: "group3",
	}

	err = synchronize(ctx, db, sc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc, expected) {
		t.Errorf(`!reflect.DeepEqual(sc, expected): %#v`, sc)
	}

	// remove a user
	db.Users = db.Users[:2]
	db.DeletedUsers = []string{"user3"}
	delete(expected.Users, "user3")

	err = synchronize(ctx, db, sc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc, expected) {
		t.Errorf(`!reflect.DeepEqual(sc, expected): %#v`, sc)
	}

	// recreate a user
	db.Users[1] = &User{
		Name:        "user2",
		UID:         2005,
		DisplayName: "u ser 2",
		Group:       "group1",
		Shell:       "/bin/zsh",
	}
	expected.Users["user2"] = &User{
		Name:        "user2",
		UID:         2005,
		DisplayName: "u ser 2",
		Group:       "group1",
		Shell:       "/bin/zsh",
	}

	err = synchronize(ctx, db, sc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc, expected) {
		t.Errorf(`!reflect.DeepEqual(sc, expected): %#v`, sc)
	}

	// update user attributes
	db.Users[1] = &User{
		Name:        "user2",
		UID:         2005,
		DisplayName: "The One",
		Group:       "group3",
		Groups:      []string{"system-group2"},
		Shell:       "/bin/bash",
		PubKeys:     []string{"pubkey1", "pubkey2"},
	}
	expected.Users["user2"] = &User{
		Name:        "user2",
		UID:         2005,
		DisplayName: "The One",
		Group:       "group3",
		Groups:      []string{"system-group2"},
		Shell:       "/bin/bash",
		PubKeys:     []string{"pubkey1", "pubkey2"},
	}

	err = synchronize(ctx, db, sc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc, expected) {
		t.Errorf(`!reflect.DeepEqual(sc, expected): %#v`, sc)
	}

	// add a group
	db.Groups = append(db.Groups, Group{Name: "group4", GID: 3003})
	expected.Groups["group4"] = &Group{Name: "group4", GID: 3003}

	err = synchronize(ctx, db, sc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc, expected) {
		t.Errorf(`!reflect.DeepEqual(sc, expected): %#v`, sc)
	}

	// remove a group
	db.Groups = db.Groups[0 : len(db.Groups)-1]
	db.DeletedGroups = append(db.DeletedGroups, "group4")
	delete(expected.Groups, "group4")

	err = synchronize(ctx, db, sc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc, expected) {
		t.Errorf(`!reflect.DeepEqual(sc, expected): %#v`, sc)
	}

	// recreate a group
	db.Groups[0] = Group{Name: "group1", GID: 3004}
	expected.Groups["group1"] = &Group{Name: "group1", GID: 3004}

	err = synchronize(ctx, db, sc)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(sc, expected) {
		t.Errorf(`!reflect.DeepEqual(sc, expected): %#v`, sc)
	}
}

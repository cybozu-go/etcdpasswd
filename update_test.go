package etcdpasswd

import (
	"context"
	"path"
	"testing"

	"github.com/coreos/etcd/clientv3"
	"github.com/cybozu-go/etcdpasswd/mocks"
	"github.com/golang/mock/gomock"
)

func clearEtcd(ctx context.Context, client Client, t *testing.T) {
	for _, key := range []string{KeyUsers, KeyDeletedUsers, KeyGroups, KeyDeletedGroups, KeyLocked} {
		_, err := client.Delete(ctx, key, clientv3.WithPrefix())
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testDeleteUser(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	clearEtcd(ctx, client, t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuth(ctrl)
	setupAuth(mockAuth)

	mockAuth.EXPECT().DeleteUser(ctx, "foo").Times(1)

	_, err := client.Put(ctx, path.Join(KeyDeletedUsers, "foo"), "")
	if err != nil {
		t.Fatal(err)
	}

	err = client.updateAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	t.Run("DeleteUser", testDeleteUser)
}

package syncer

import (
	"os/user"
	"testing"
)

func TestMakeUser(t *testing.T) {
	t.Parallel()

	uu, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	u, err := makeUser(uu)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", u)
}

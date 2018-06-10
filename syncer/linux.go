package syncer

import (
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"github.com/cybozu-go/etcdpasswd"
)

func gid2Name(gid string) string {
	g, err := user.LookupGroupId(gid)
	if err != nil {
		return gid
	}
	return g.Name
}

func makeUser(uu *user.User) (*etcdpasswd.User, error) {
	u := &etcdpasswd.User{
		Name:        uu.Username,
		DisplayName: uu.Name,
	}

	uid, err := strconv.Atoi(uu.Uid)
	if err != nil {
		return nil, err
	}
	u.UID = uid
	u.Group = gid2Name(uu.Gid)

	gids, err := uu.GroupIds()
	if err != nil {
		return nil, err
	}
	gnames := make([]string, len(gids))
	for i, gid := range gids {
		gnames[i] = gid2Name(gid)
	}
	u.Groups = gnames

	out, err := exec.Command("getent", "passwd", uu.Name).Output()
	if err != nil {
		return nil, err
	}
	fields := strings.Split(strings.TrimSpace(string(out)), ":")
	u.Shell = fields[len(fields)-1]

	pubkeys, err := getPubKeys(uu.HomeDir)
	if err != nil {
		return nil, err
	}
	u.PubKeys = pubkeys

	return u, nil
}

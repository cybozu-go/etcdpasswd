package syncer

import (
	"context"
	"os/user"
	"strconv"
	"strings"

	"github.com/cybozu-go/etcdpasswd"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
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
	gnames := make([]string, 0, len(gids))
	for _, gid := range gids {
		if gid == uu.Gid {
			continue
		}
		gnames = append(gnames, gid2Name(gid))
	}
	u.Groups = gnames

	c := well.CommandContext(context.Background(), "getent", "passwd", uu.Username)
	c.Severity = log.LvDebug
	out, err := c.Output()
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

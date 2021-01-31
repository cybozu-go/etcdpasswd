package mtest

import (
	"strconv"
	"time"
)

type generator struct {
	userPrefix  string
	groupPrefix string
	defaultUID  int64
	defaultGID  int64
	userCount   int
	groupCount  int
}

var gen = newGenerator()

func newGenerator() *generator {
	now := time.Now()

	user := "epu" + strconv.FormatInt(now.Unix(), 10)
	group := "epg" + strconv.FormatInt(now.Unix(), 10)
	defaultUID := now.Unix()
	defaultGID := now.Unix()

	return &generator{
		userPrefix:  user,
		groupPrefix: group,
		defaultUID:  defaultUID,
		defaultGID:  defaultGID,
	}
}

func (g *generator) newUsername() string {
	g.userCount++
	return g.userPrefix + strconv.Itoa(g.userCount)
}

func (g *generator) newGroupname() string {
	g.groupCount++
	return g.groupPrefix + strconv.Itoa(g.groupCount)
}

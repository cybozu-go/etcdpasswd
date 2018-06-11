package cli

import (
	"strings"

	"github.com/cybozu-go/etcdpasswd"
)

var client etcdpasswd.Client

// Setup setups this package.
func Setup(c etcdpasswd.Client) {
	client = c
}

type commaStrings []string

func (o *commaStrings) String() string {
	return strings.Join([]string(*o), ",")
}

func (o *commaStrings) Set(v string) error {
	*o = strings.Split(v, ",")
	return nil
}

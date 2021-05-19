package cmd

import (
	"strings"
)

type commaStrings []string

func (o *commaStrings) String() string {
	return strings.Join([]string(*o), ",")
}

func (o *commaStrings) Set(v string) error {
	*o = strings.Split(v, ",")
	return nil
}

func (o *commaStrings) Type() string {
	return "commaStrings"
}

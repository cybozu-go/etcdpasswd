package itest

import (
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	bridgeAddress = os.Getenv("BRIDGE_ADDRESS")
	host1         = os.Getenv("HOST1")
	host2         = os.Getenv("HOST2")
	host3         = os.Getenv("HOST3")
	placemat      = os.Getenv("PLACEMAT")
	sshKey        ssh.Signer
	debug         = os.Getenv("DEBUG") == "1"
)

func init() {
	s, err := os.Open(os.Getenv("SSH_PRIVKEY"))
	if err != nil {
		panic(err)
	}
	defer s.Close()

	data, err := ioutil.ReadAll(s)
	if err != nil {
		panic(err)
	}

	sshKey, err = ssh.ParsePrivateKey(data)
	if err != nil {
		panic(err)
	}
}

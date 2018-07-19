package itest

import (
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
)

func runPlacemat() *exec.Cmd {
	cmd := exec.Command(placemat, "cluster.yml")
	if debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	return cmd
}

func sshTo(address string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(sshKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	return ssh.Dial("tcp", address+":22", config)
}

func prepareSshClients(addresses ...string) map[string]*ssh.Client {
	ret := make(map[string]*ssh.Client)
	ch := time.After(time.Minute)
	for _, a := range addresses {
	RETRY:
		select {
		case <-ch:
			panic("timed out")
		default:
		}
		client, err := sshTo(a)
		if err != nil {
			time.Sleep(time.Second)
			goto RETRY
		}
		ret[a] = client
	}
	return ret
}

func runEtcd(client *ssh.Client) error {
	return nil
}

package itest

import (
	"bytes"
	"errors"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	sshClients = make(map[string]*ssh.Client)
)

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

func prepareSshClients(addresses ...string) error {
	ch := time.After(time.Minute)
	for _, a := range addresses {
	RETRY:
		select {
		case <-ch:
			return errors.New("timed out")
		default:
		}
		client, err := sshTo(a)
		if err != nil {
			time.Sleep(time.Second)
			goto RETRY
		}
		sshClients[a] = client
	}
	return nil
}

func runEtcd(client *ssh.Client) error {
	command := "systemd-run --user /data/etcd --listen-client-urls=http://0.0.0.0:2379 --advertise-client-urls=http://localhost:2379"
	sess, err := client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	return sess.Run(command)
}

func runEPAgent() error {
	for _, c := range sshClients {
		sess, err := c.NewSession()
		if err != nil {
			return err
		}

		err = sess.Run("sudo systemd-run /data/ep-agent")
		sess.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func runCommand(host, command string) (stdout, stderr []byte, e error) {
	client := sshClients[host]
	sess, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	defer sess.Close()

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	sess.Stdout = outBuf
	sess.Stderr = errBuf
	err = sess.Run(command)
	return outBuf.Bytes(), errBuf.Bytes(), err
}

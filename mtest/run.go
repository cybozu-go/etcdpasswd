package mtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cybozu-go/well"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
)

const (
	sshTimeout         = 3 * time.Minute
	defaultDialTimeout = 30 * time.Second
	defaultKeepAlive   = 5 * time.Second

	// DefaultRunTimeout is the timeout value for Agent.Run().
	DefaultRunTimeout = 10 * time.Minute
)

var (
	sshClients = make(map[string]*sshAgent)
	httpClient = &well.HTTPClient{Client: &http.Client{}}

	agentDialer = &net.Dialer{
		Timeout:   defaultDialTimeout,
		KeepAlive: defaultKeepAlive,
	}
)

type sshAgent struct {
	client *ssh.Client
	conn   net.Conn
}

func sshTo(address string, sshKey ssh.Signer, userName string) (*sshAgent, error) {
	conn, err := agentDialer.Dial("tcp", address+":22")
	if err != nil {
		fmt.Printf("failed to dial: %s\n", address)
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(sshKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	err = conn.SetDeadline(time.Now().Add(defaultDialTimeout))
	if err != nil {
		conn.Close()
		return nil, err
	}
	clientConn, channelCh, reqCh, err := ssh.NewClientConn(conn, "tcp", config)
	if err != nil {
		// conn was already closed in ssh.NewClientConn
		return nil, err
	}
	err = conn.SetDeadline(time.Time{})
	if err != nil {
		clientConn.Close()
		return nil, err
	}
	a := sshAgent{
		client: ssh.NewClient(clientConn, channelCh, reqCh),
		conn:   conn,
	}
	return &a, nil
}

func prepareSSHClients(addresses ...string) error {
	sshKey, err := parsePrivateKey(sshKeyFile)
	if err != nil {
		return err
	}

	ch := time.After(sshTimeout)
	for _, a := range addresses {
	RETRY:
		select {
		case <-ch:
			return errors.New("prepareSSHClients timed out")
		default:
		}
		agent, err := sshTo(a, sshKey, "ubuntu")
		if err != nil {
			time.Sleep(time.Second)
			goto RETRY
		}
		sshClients[a] = agent
	}

	return nil
}

func parsePrivateKey(keyPath string) (ssh.Signer, error) {
	f, err := os.Open(keyPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(data)
}

func stopEtcd() error {
	env := well.NewEnvironment(context.Background())
	for _, host := range []string{host1, host2, host3} {
		host2 := host
		env.Go(func(ctx context.Context) error {
			sess, err := sshClients[host2].client.NewSession()
			if err != nil {
				return err
			}
			defer sess.Close()
			return sess.Run("sudo systemctl stop my-etcd.service; sudo rm -rf /home/ubuntu/default.etcd")
		})
	}
	env.Stop()
	return env.Wait()
}

func runEtcd() error {
	env := well.NewEnvironment(context.Background())
	for _, host := range []string{host1, host2, host3} {
		host2 := host
		env.Go(func(ctx context.Context) error {
			sess, err := sshClients[host2].client.NewSession()
			if err != nil {
				return err
			}
			defer sess.Close()
			return sess.Run("sudo systemd-run --unit=my-etcd.service /opt/bin/etcd --listen-client-urls=http://0.0.0.0:2379 --advertise-client-urls=http://localhost:2379 --data-dir /home/ubuntu/default.etcd")
		})
	}
	env.Stop()
	return env.Wait()
}

func stopEPAgent() error {
	env := well.NewEnvironment(context.Background())
	for _, host := range []string{host1, host2, host3} {
		host2 := host
		env.Go(func(ctx context.Context) error {
			sess, err := sshClients[host2].client.NewSession()
			if err != nil {
				return err
			}
			defer sess.Close()
			sess.Run("sudo systemctl reset-failed ep-agent.service; sudo systemctl stop ep-agent.service")
			return nil // Ignore error if ep-agent was not running
		})
	}
	env.Stop()
	return env.Wait()
}

func runEPAgent() error {
	env := well.NewEnvironment(context.Background())
	for _, host := range []string{host1, host2, host3} {
		host2 := host
		env.Go(func(ctx context.Context) error {
			sess, err := sshClients[host2].client.NewSession()
			if err != nil {
				return err
			}
			defer sess.Close()
			return sess.Run("sudo systemd-run --unit=ep-agent.service /opt/bin/ep-agent")
		})
	}
	env.Stop()
	return env.Wait()
}

func execAt(host string, args ...string) (stdout, stderr []byte, e error) {
	return execAtWithStream(host, nil, args...)
}

// WARNING: `input` can contain secret data.  Never output `input` to console.
func execAtWithInput(host string, input []byte, args ...string) (stdout, stderr []byte, e error) {
	var r io.Reader
	if input != nil {
		r = bytes.NewReader(input)
	}
	return execAtWithStream(host, r, args...)
}

// WARNING: `input` can contain secret data.  Never output `input` to console.
func execAtWithStream(host string, input io.Reader, args ...string) (stdout, stderr []byte, e error) {
	agent := sshClients[host]
	return doExec(agent, input, args...)
}

// WARNING: `input` can contain secret data.  Never output `input` to console.
func doExec(agent *sshAgent, input io.Reader, args ...string) ([]byte, []byte, error) {
	err := agent.conn.SetDeadline(time.Now().Add(DefaultRunTimeout))
	if err != nil {
		return nil, nil, err
	}
	defer agent.conn.SetDeadline(time.Time{})

	sess, err := agent.client.NewSession()
	if err != nil {
		return nil, nil, err
	}
	defer sess.Close()

	if input != nil {
		sess.Stdin = input
	}
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	sess.Stdout = outBuf
	sess.Stderr = errBuf
	err = sess.Run(strings.Join(args, " "))
	return outBuf.Bytes(), errBuf.Bytes(), err
}

func execSafeAt(host string, args ...string) []byte {
	stdout, stderr, err := execAt(host, args...)
	ExpectWithOffset(1, err).To(Succeed(), "[%s] %v: %s", host, args, stderr)
	return stdout
}

func etcdpasswd(args ...string) ([]byte, []byte, error) {
	args = append([]string{"/opt/bin/etcdpasswd"}, args...)
	return execAt(host1, args...)
}

func etcdpasswdSafe(args ...string) []byte {
	args = append([]string{"/opt/bin/etcdpasswd"}, args...)
	return execSafeAt(host1, args...)
}

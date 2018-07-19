package itest

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestItest(t *testing.T) {
	if len(sshKeyFile) == 0 {
		t.Skip("no SSH_PRIVKEY envvar")
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration test for etcdpasswd")
}

var _ = BeforeSuite(func() {
	fmt.Println("Preparing...")
	err := prepareSshClients(host1, host2, host3)
	Expect(err).NotTo(HaveOccurred())

	// sync VM root filesystem to store newly generated SSH host keys.
	for h := range sshClients {
		execSafeAt(h, "sync")
	}

	err = runEtcd(sshClients[host1])
	Expect(err).NotTo(HaveOccurred())

	time.Sleep(time.Second)

	err = runEPAgent()
	Expect(err).NotTo(HaveOccurred())

	time.Sleep(time.Second)
	fmt.Println("Begin tests...")
})

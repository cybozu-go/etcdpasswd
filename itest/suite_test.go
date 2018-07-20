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

	SetDefaultEventuallyPollingInterval(time.Second)
	SetDefaultEventuallyTimeout(time.Minute)

	err := prepareSSHClients(host1, host2, host3)
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

	By("get start-uid should return 0")
	stdout := execSafeAt(host1, CLI, "get", "start-uid")
	Expect(stdout).To(Equal("0\n"))

	By("set start-uid should succeed")
	execSafeAt(host1, CLI, "set", "start-uid", "2000")
	stdout = execSafeAt(host1, CLI, "get", "start-uid")
	Expect(stdout).To(Equal("2000\n"))

	By("get start-gid should return 0")
	stdout = execSafeAt(host1, CLI, "get", "start-gid")
	Expect(stdout).To(Equal("0\n"))

	By("set start-gid should succeed")
	execSafeAt(host1, CLI, "set", "start-gid", "2000")
	stdout = execSafeAt(host1, CLI, "get", "start-gid")
	Expect(stdout).To(Equal("2000\n"))

	By("get default-group should be empty")
	stdout = execSafeAt(host1, CLI, "get", "default-group")
	Expect(stdout).To(Equal("\n"))

	By("set default-group should succeed")
	execSafeAt(host1, CLI, "set", "default-group", "cybozu")
	stdout = execSafeAt(host1, CLI, "get", "default-group")
	Expect(stdout).To(Equal("cybozu\n"))

	By("get default-groups should be empty")
	stdout = execSafeAt(host1, CLI, "get", "default-groups")
	Expect(stdout).To(BeEmpty())

	By("set default-groups should succeed")
	execSafeAt(host1, CLI, "set", "default-groups", "sudo,adm")
	stdout = execSafeAt(host1, CLI, "get", "default-groups")
	Expect(stdout).To(MatchRegexp("\\bsudo\\b"))
	Expect(stdout).To(MatchRegexp("\\badm\\b"))

})

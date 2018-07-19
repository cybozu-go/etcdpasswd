package itest

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestItest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Itest Suite")
}

var _ = BeforeSuite(func() {
	err := prepareSshClients(host1, host2, host3)
	Expect(err).NotTo(HaveOccurred())

	err = runEtcd(sshClients[host1])
	Expect(err).NotTo(HaveOccurred())

	time.Sleep(time.Second)

	err = runEPAgent()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("etcdpasswd is working", func() {
	It("should not fail", func(done Done) {
		_, _, err := runCommand(host1, "/data/etcdpasswd get start-uid")
		Expect(err).NotTo(HaveOccurred())
		close(done)
	}, 10)
})

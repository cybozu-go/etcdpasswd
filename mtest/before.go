package mtest

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cybozu-go/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// RunBeforeSuite is for Ginkgo BeforeSuite
func RunBeforeSuite() {
	fmt.Println("Preparing...")

	SetDefaultEventuallyPollingInterval(time.Second)
	SetDefaultEventuallyTimeout(time.Minute)

	log.DefaultLogger().SetThreshold(log.LvError)

	err := prepareSSHClients(host1, host2, host3)
	Expect(err).NotTo(HaveOccurred())

	// sync VM root filesystem to store newly generated SSH host keys.
	for h := range sshClients {
		execSafeAt(h, "sync")
	}

	By("copying test files")
	for _, testFile := range []string{etcdPath, etcdctlPath, etcdpasswdPath, epagentPath} {
		f, err := os.Open(testFile)
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()
		remoteFilename := filepath.Join("/tmp", filepath.Base(testFile))
		for _, host := range []string{host1, host2, host3} {
			_, err := f.Seek(0, os.SEEK_SET)
			Expect(err).NotTo(HaveOccurred())
			stdout, stderr, err := execAt(host, "sudo", "mkdir", "-p", "/opt/bin")
			Expect(err).NotTo(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)
			stdout, stderr, err = execAtWithStream(host, f, "dd", "of="+remoteFilename)
			Expect(err).NotTo(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)
			stdout, stderr, err = execAt(host, "sudo", "mv", remoteFilename, filepath.Join("/opt/bin", filepath.Base(testFile)))
			Expect(err).NotTo(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)
			stdout, stderr, err = execAt(host, "sudo", "chmod", "755", filepath.Join("/opt/bin", filepath.Base(testFile)))
			Expect(err).NotTo(HaveOccurred(), "stdout=%s, stderr=%s", stdout, stderr)
		}
	}

	By("starting etcd")
	err = stopEtcd()
	Expect(err).NotTo(HaveOccurred())
	err = runEtcd()
	Expect(err).NotTo(HaveOccurred())

	time.Sleep(time.Second)

	By("starting ep-agent")
	err = stopEPAgent()
	Expect(err).NotTo(HaveOccurred())
	err = runEPAgent()
	Expect(err).NotTo(HaveOccurred())

	time.Sleep(time.Second)
	fmt.Println("Begin tests...")

	By("get start-uid should return 0")
	stdout, stderr, err := etcdpasswd("get", "start-uid")
	Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
	Expect(string(stdout)).To(Equal("0\n"))

	By("set start-uid should succeed")
	defaultUID := strconv.FormatInt(gen.defaultUID, 10)
	etcdpasswdSafe("set", "start-uid", defaultUID)
	stdout, stderr, err = etcdpasswd("get", "start-uid")
	Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
	Expect(strings.TrimSpace(string(stdout))).To(Equal(defaultUID))

	By("get start-gid should return 0")
	stdout, stderr, err = etcdpasswd("get", "start-gid")
	Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
	Expect(string(stdout)).To(Equal("0\n"))

	By("set start-gid should succeed")
	defaultGID := strconv.FormatInt(gen.defaultGID, 10)
	etcdpasswdSafe("set", "start-gid", defaultGID)
	stdout, stderr, err = etcdpasswd("get", "start-gid")
	Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
	Expect(strings.TrimSpace(string(stdout))).To(Equal(defaultGID))

	By("get default-group should be empty")
	stdout, stderr, err = etcdpasswd("get", "default-group")
	Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
	Expect(string(stdout)).To(Equal("\n"))

	By("set default-group should succeed")
	etcdpasswdSafe("set", "default-group", "users")
	stdout, stderr, err = etcdpasswd("get", "default-group")
	Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
	Expect(string(stdout)).To(Equal("users\n"))

	By("get default-groups should be empty")
	stdout, stderr, err = etcdpasswd("get", "default-groups")
	Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
	Expect(string(stdout)).To(BeEmpty())

	By("set default-groups should succeed")
	etcdpasswdSafe("set", "default-groups", "sudo,adm")
	stdout, stderr, err = etcdpasswd("get", "default-groups")
	Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
	Expect(string(stdout)).To(MatchRegexp("\\bsudo\\b"))
	Expect(string(stdout)).To(MatchRegexp("\\badm\\b"))
}

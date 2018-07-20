package itest

import (
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// CLI is the file path of etcdpasswd
const CLI = "/data/etcdpasswd"

type testCounter struct {
	n   int
	mux sync.Mutex
}

func (c *testCounter) next() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.n++
	return c.n
}

var counter testCounter

type testContext struct {
	user  string
	group string
}

const (
	testUserBase  = "epu"
	testGroupBase = "epg"
)

func setupTest() *testContext {
	n := counter.next()
	return &testContext{
		user:  testUserBase + strconv.Itoa(n),
		group: testGroupBase + strconv.Itoa(n),
	}
}

var _ = Describe("etcdpasswd", func() {
	hosts := []string{host1, host2, host3}

	It("group add/remove, user add/remove", func() {
		c := setupTest()

		By("create group and node users")
		execSafeAt(host1, CLI, "group", "add", c.group)
		execSafeAt(host1, CLI, "user", "add", "-group", c.group, c.user)

		By("should create group and user")
		stdout := execSafeAt(host1, CLI, "user", "list")
		Expect(stdout).To(MatchRegexp("\\b%s\\b", c.user))

		stdout = execSafeAt(host1, CLI, "group", "list")
		Expect(stdout).To(MatchRegexp("\\b%s\\b", c.group))

		for _, h := range hosts {
			Eventually(func() int {
				stdout, _, err := execAt(h, "id", "-u", c.user)
				if err != nil {
					return -1
				}
				uid, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
				if err != nil {
					return -1
				}
				return uid
			}).Should(BeNumerically(">=", 2000))

			Eventually(func() int {
				stdout, _, err := execAt(h, "id", "-g", c.user)
				if err != nil {
					return -1
				}
				gid, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
				if err != nil {
					return -1
				}
				return gid
			}).Should(BeNumerically(">=", 2000))

			Eventually(func() []string {
				stdout, _, err := execAt(h, "id", "-Gn", c.user)
				if err != nil {
					return []string{}
				}
				groups := strings.Split(strings.TrimSpace(string(stdout)), " ")
				return groups
			}).Should(ConsistOf(c.group, "sudo", "adm"))
		}

		By("should remove user and group")
		execAt(host1, CLI, "user", "remove", c.user)
		execAt(host1, CLI, "group", "remove", c.group)

		stdout = execSafeAt(host1, CLI, "user", "list")
		Expect(stdout).NotTo(MatchRegexp("\\b%s\\b", c.user))

		stdout = execSafeAt(host1, CLI, "group", "list")
		Expect(stdout).NotTo(MatchRegexp("\\b%s\\b", c.group))

		for _, h := range hosts {
			Eventually(func() error {
				_, _, err := execAt(h, "id", c.user)
				return err
			}).ShouldNot(Exit(0))

			Eventually(func() error {
				_, _, err := execAt(h, "getent", "group", c.group)
				return err
			}).ShouldNot(Exit(0))
		}
	})

	Context("cert add", func() {
		c := setupTest()

		var dir string

		BeforeEach(func() {
			stdout := execSafeAt(host1, "mktemp", "-d")
			dir = strings.TrimSpace(stdout)
			execSafeAt(host1, CLI, "user", "add", c.user)
			execSafeAt(host1, "ssh-keygen", "-t", "rsa", "-N", "''", "-C", "'test cert add'", "-f", filepath.Join(dir, "id_rsa"))
			execSafeAt(host1, CLI, "cert", "add", c.user, filepath.Join(dir, "id_rsa.pub"))
		})

		AfterEach(func() {
			execAt(host1, "rm", "-rf", dir)
			execAt(host1, CLI, "cert", "remove", c.user, "0")
			execAt(host1, CLI, "user", "remove", c.user)
		})

		It("should add SSH key", func() {
			stdout := execSafeAt(host1, CLI, "cert", "list", c.user)
			Expect(stdout).To(ContainSubstring("test cert add"))

			for _, h := range hosts {
				stdout = execSafeAt(h, "sudo", "ssh-keygen", "-l", "-f", "/home/"+c.user+"/.ssh/authorized_keys")
				Expect(stdout).To(ContainSubstring("test cert add"))
			}
		})
	})

	Context("cert remove", func() {
		c := setupTest()
		var dir string

		BeforeEach(func() {
			stdout := execSafeAt(host1, "mktemp", "-d")
			dir = strings.TrimSpace(stdout)
			execSafeAt(host1, CLI, "user", "add", c.user)
			execSafeAt(host1, "ssh-keygen", "-t", "rsa", "-N", "''", "-C", "'test cert remove'", "-f", filepath.Join(dir, "id_rsa"))
			execSafeAt(host1, CLI, "cert", "add", c.user, filepath.Join(dir, "id_rsa.pub"))
		})

		AfterEach(func() {
			execAt(host1, "rm", "-rf", dir)
			execAt(host1, CLI, "user", "remove", c.user)
		})

		It("should remove SSH key", func() {
			execSafeAt(host1, CLI, "cert", "remove", c.user, "0")

			stdout := execSafeAt(host1, CLI, "cert", "list", c.user)
			Expect(stdout).To(BeZero())

			for _, h := range hosts {
				stdout = execSafeAt(h, "sudo", "cat", "/home/"+c.user+"/.ssh/authorized_keys")
				keys := strings.TrimSpace(stdout)
				Expect(keys).To(BeZero())
			}
		})
	})

	Context("locker add", func() {
		c := setupTest()

		BeforeEach(func() {
			execSafeAt(host1, CLI, "user", "add", c.user)
			// Created user is locked by default.
			// "usermod -U" requires that non-empty password be set.
			// "usermod -p" takes *encrypted* password as its argument,
			// so string "invalid" does not result a weak password.
			for _, h := range hosts {
				execSafeAt(h, "sudo", "usermod", "-p", "invalid", c.user)
				execSafeAt(h, "sudo", "usermod", "-U", c.user)
			}
		})

		AfterEach(func() {
			execAt(host1, CLI, "user", "remove", c.user)
		})

		It("should lock user", func() {
			execSafeAt(host1, CLI, "locker", "add", c.user)

			stdout := execSafeAt(host1, CLI, "locker", "list")
			Expect(stdout).To(MatchRegexp("\\b%s\\b", c.user))

			for _, h := range hosts {
				stdout = execSafeAt(h, "sudo", "grep", "^"+c.user+":!", "/etc/shadow")
				Expect(stdout).NotTo(BeZero())
			}
		})
	})

	Context("locker remove", func() {
		c := setupTest()

		BeforeEach(func() {
			execSafeAt(host1, CLI, "user", "add", c.user)
			execSafeAt(host1, CLI, "locker", "add", c.user)
		})

		AfterEach(func() {
			execAt(host1, CLI, "user", "remove", c.user)
		})

		It("should remove user from list", func() {
			execSafeAt(host1, CLI, "locker", "remove", c.user)

			stdout := execSafeAt(host1, CLI, "locker", "list")
			Expect(stdout).NotTo(MatchRegexp("\\b%s\\b", c.user))
		})
	})
})

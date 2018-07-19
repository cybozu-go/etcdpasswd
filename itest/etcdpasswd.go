package itest

import (
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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

	Context("etcdpasswd get start-uid", func() {
		It("should return 2000", func() {
			stdout := execSafeAt(host1, CLI, "get", "start-uid")
			Expect(stdout).To(Equal("2000\n"))
		})
	})

	Context("etcdpasswd get start-gid", func() {
		It("should return 2000", func() {
			stdout := execSafeAt(host1, CLI, "get", "start-gid")
			Expect(stdout).To(Equal("2000\n"))
		})
	})

	Context("etcdpasswd get default-group", func() {
		It("should return cybozu", func() {
			stdout := execSafeAt(host1, CLI, "get", "default-group")
			Expect(stdout).To(Equal("cybozu\n"))
		})
	})

	Context("etcdpasswd get default-groups", func() {
		It("should return sudo and adm", func() {
			stdout := execSafeAt(host1, CLI, "get", "default-groups")
			Expect(stdout).To(MatchRegexp("\\bsudo\\b"))
			Expect(stdout).To(MatchRegexp("\\badm\\b"))
		})
	})

	Context("group add, user add", func() {
		c := setupTest()

		BeforeEach(func() {
			execSafeAt(host1, CLI, "group", "add", c.group)
			execSafeAt(host1, CLI, "user", "add", "-group", c.group, c.user)
		})

		AfterEach(func() {
			execAt(host1, CLI, "user", "remove", c.user)
			execAt(host1, CLI, "group", "remove", c.group)
		})

		It("should create group and user", func() {
			stdout := execSafeAt(host1, CLI, "user", "list")
			Expect(stdout).To(MatchRegexp("\\b%s\\b", c.user))

			stdout = execSafeAt(host1, CLI, "group", "list")
			Expect(stdout).To(MatchRegexp("\\b%s\\b", c.group))

			for _, h := range hosts {
				stdout = execSafeAt(h, "id", "-u", c.user)
				uid, err := strconv.Atoi(strings.TrimSpace(stdout))
				if err != nil {
					Fail(err.Error())
				}
				Expect(uid).To(BeNumerically(">=", 2000))

				stdout = execSafeAt(h, "id", "-g", c.user)
				gid, err := strconv.Atoi(strings.TrimSpace(stdout))
				if err != nil {
					Fail(err.Error())
				}
				Expect(gid).To(BeNumerically(">=", 2000))

				stdout = execSafeAt(h, "id", "-Gn", c.user)
				groups := strings.Split(strings.TrimSpace(stdout), " ")
				Expect(groups).To(ConsistOf(c.group, "sudo", "adm"))
			}
		})
	})

	Context("remove user, remove group", func() {
		c := setupTest()

		BeforeEach(func() {
			execSafeAt(host1, CLI, "group", "add", c.group)
			execSafeAt(host1, CLI, "user", "add", "-group", c.group, c.user)
		})

		It("should remove user and group", func() {
			execAt(host1, CLI, "user", "remove", c.user)
			execAt(host1, CLI, "group", "remove", c.group)

			stdout := execSafeAt(host1, CLI, "user", "list")
			Expect(stdout).NotTo(MatchRegexp("\\b%s\\b", c.user))

			stdout = execSafeAt(host1, CLI, "group", "list")
			Expect(stdout).NotTo(MatchRegexp("\\b%s\\b", c.group))

			for _, h := range hosts {
				_, _, err := execAt(h, "id", c.user)
				Expect(err).NotTo(Exit(0))

				_, _, err = execAt(h, "getent", "group", c.group)
				Expect(err).NotTo(Exit(0))
			}
		})
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

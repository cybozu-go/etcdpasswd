package mtest

import (
	"path/filepath"
	"strconv"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testEtcdpasswd() {
	hosts := []string{host1, host2, host3}

	It("group add/remove, user add/remove", func() {
		group := gen.newGroupname()
		user := gen.newUsername()

		By("create group and node users")
		etcdpasswdSafe("group", "add", group)
		etcdpasswdSafe("user", "add", "-group", group, user)

		By("should create group and user")
		stdout, stderr, err := etcdpasswd("user", "list")
		Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
		Expect(stdout).To(MatchRegexp("\\b%s\\b", user))

		stdout, stderr, err = etcdpasswd("group", "list")
		Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
		Expect(stdout).To(MatchRegexp("\\b%s\\b", group))

		for _, h := range hosts {
			Eventually(func() int {
				stdout, _, err := execAt(h, "id", "-u", user)
				if err != nil {
					return -1
				}
				uid, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
				if err != nil {
					return -1
				}
				return uid
			}).Should(BeNumerically(">=", gen.defaultUID))

			Eventually(func() int {
				stdout, _, err := execAt(h, "id", "-g", user)
				if err != nil {
					return -1
				}
				gid, err := strconv.Atoi(strings.TrimSpace(string(stdout)))
				if err != nil {
					return -1
				}
				return gid
			}).Should(BeNumerically(">=", gen.defaultGID))

			Eventually(func() []string {
				stdout, _, err := execAt(h, "id", "-Gn", user)
				if err != nil {
					return []string{}
				}
				groups := strings.Split(strings.TrimSpace(string(stdout)), " ")
				return groups
			}).Should(ConsistOf(group, "sudo", "adm"))
		}

		By("should remove user and group")
		etcdpasswdSafe("user", "remove", user)
		etcdpasswdSafe("group", "remove", group)

		stdout, stderr, err = etcdpasswd("user", "list")
		Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
		Expect(stdout).NotTo(MatchRegexp("\\b%s\\b", user))

		stdout, stderr, err = etcdpasswd("group", "list")
		Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
		Expect(stdout).NotTo(MatchRegexp("\\b%s\\b", group))

		for _, h := range hosts {
			Eventually(func() error {
				_, _, err := execAt(h, "id", user)
				return err
			}).ShouldNot(Exit(0))

			Eventually(func() error {
				_, _, err := execAt(h, "getent", "group", group)
				return err
			}).ShouldNot(Exit(0))
		}
	})

	It("cert add/remove", func() {
		user := gen.newUsername()

		By("Create user and ssh key")
		stdout := execSafeAt(host1, "mktemp", "-d")
		dir := strings.TrimSpace(string(stdout))
		etcdpasswdSafe("user", "add", user)
		execSafeAt(host1, "ssh-keygen", "-t", "rsa", "-N", "''", "-C", "'test cert add'", "-f", filepath.Join(dir, "id_rsa"))
		etcdpasswdSafe("cert", "add", user, filepath.Join(dir, "id_rsa.pub"))

		By("should add SSH key")
		stdout, stderr, err := etcdpasswd("cert", "list", user)
		Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
		Expect(stdout).To(ContainSubstring("test cert add"))

		for _, h := range hosts {
			Eventually(func() string {
				stdout, _, err := execAt(h, "sudo", "ssh-keygen", "-l", "-f", "/home/"+user+"/.ssh/authorized_keys")
				if err != nil {
					return ""
				}
				return string(stdout)
			}).Should(ContainSubstring("test cert add"))
		}

		By("should remove SSH key")
		etcdpasswdSafe("cert", "remove", user, "0")

		stdout, stderr, err = etcdpasswd("cert", "list", user)
		Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
		Expect(string(stdout)).To(BeZero())

		for _, h := range hosts {
			Eventually(func() string {
				stdout, _, err := execAt(h, "sudo", "ssh-keygen", "-l", "-f", "/home/"+user+"/.ssh/authorized_keys")
				if err == nil {
					return "test cert add"
				}
				return string(stdout)
			}).ShouldNot(ContainSubstring("test cert add"))
		}
	})

	It("locker add/remove", func() {
		user := gen.newUsername()

		By("user add")
		etcdpasswdSafe("user", "add", user)
		// Created user is locked by default.
		// "usermod -U" requires that non-empty password be set.
		// "usermod -p" takes *encrypted* password as its argument,
		// so string "invalid" does not result a weak password.
		for _, h := range hosts {
			Eventually(func() error {
				_, _, err := execAt(h, "sudo", "usermod", "-p", "invalid", user)
				return err
			}).Should(Succeed())

			Eventually(func() error {
				_, _, err := execAt(h, "sudo", "usermod", "-U", user)
				return err
			}).Should(Succeed())
		}

		By("should lock user")
		etcdpasswdSafe("locker", "add", user)

		stdout, stderr, err := etcdpasswd("locker", "list")
		Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
		Expect(string(stdout)).To(MatchRegexp("\\b%s\\b", user))

		for _, h := range hosts {
			Eventually(func() string {
				stdout, _, err := execAt(h, "sudo", "grep", "^"+user+":!", "/etc/shadow")
				if err != nil {
					return ""
				}
				return string(stdout)
			}).ShouldNot(BeZero())
		}

		By("should remove user from list")
		etcdpasswdSafe("locker", "remove", user)

		stdout, stderr, err = etcdpasswd("locker", "list")
		Expect(err).NotTo(HaveOccurred(), "stderr=%s", stderr)
		Expect(string(stdout)).NotTo(MatchRegexp("\\b%s\\b", user))
	})
}

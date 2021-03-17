package syncer

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/cybozu-go/log"
)

func parseAuthorizedKeys(p string) ([]string, error) {
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var ret []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		t := scanner.Text()

		// ignore empty and comment lines
		if len(t) == 0 || t[0] == '#' {
			continue
		}
		ret = append(ret, t)
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func getPubKeys(homedir string) ([]string, error) {
	pubkeys, err := parseAuthorizedKeys(filepath.Join(homedir, ".ssh", "authorized_keys"))
	if err != nil {
		return nil, err
	}

	pubkeys2, err := parseAuthorizedKeys(filepath.Join(homedir, ".ssh", "authorized_keys2"))
	if err != nil {
		return nil, err
	}

	return append(pubkeys, pubkeys2...), nil
}

func savePubKeys(homedir string, uid, gid int, pubkeys []string) error {
	sshDir := filepath.Join(homedir, ".ssh")
	_, err := os.Stat(sshDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err = os.Mkdir(sshDir, 0700)
		if err != nil {
			return err
		}
		err = os.Chown(sshDir, uid, gid)
		if err != nil {
			return err
		}
	}

	// remove alternative file, if any.
	err = os.Remove(filepath.Join(sshDir, "authorized_keys2"))
	if err == nil {
		log.Info("removed authorized_keys2", map[string]interface{}{
			"dir": sshDir,
		})
	}

	f, err := os.CreateTemp(sshDir, ".gp")
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Chown(uid, gid)
	if err != nil {
		return err
	}
	err = f.Chmod(0600)
	if err != nil {
		return err
	}

	_, err = f.WriteString(strings.Join(pubkeys, "\n") + "\n")
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}

	return os.Rename(f.Name(), filepath.Join(sshDir, "authorized_keys"))
}

package syncer

import (
	"bufio"
	"os"
	"path/filepath"
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

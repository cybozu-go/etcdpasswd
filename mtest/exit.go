package mtest

import (
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"golang.org/x/crypto/ssh"
)

// Exit is a mathcer for exit status
func Exit(code int) types.GomegaMatcher {
	return &exitMatcher{
		exitCode: code,
	}
}

type exitMatcher struct {
	exitCode       int
	actualExitCode int
}

func (m *exitMatcher) Match(actual interface{}) (success bool, err error) {
	if m.exitCode == 0 && actual == nil {
		return true, nil
	}

	switch ee := actual.(type) {
	case *exec.ExitError:
		ws, ok := ee.Sys().(syscall.WaitStatus)
		if !ok {
			return false, errors.New("failed to obtain syscall.WaitStatus")
		}
		m.actualExitCode = ws.ExitStatus()
		return m.exitCode == m.actualExitCode, nil
	case *ssh.ExitError:
		m.actualExitCode = ee.ExitStatus()
		return m.exitCode == m.actualExitCode, nil
	default:
		return false, fmt.Errorf("invalid type: %T", actual)
	}

}

func (m *exitMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(m.actualExitCode, "to match exit code:", m.exitCode)
}

func (m *exitMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(m.actualExitCode, "not to match exit code:", m.exitCode)
}

package mtest

import . "github.com/onsi/ginkgo"

// FunctionsSuite is a test suite that tests small test cases
var FunctionsSuite = func() {
	Context("etcdpasswd", TestEtcdpasswd)
}

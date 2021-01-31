module github.com/cybozu-go/etcdpasswd

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/cybozu-go/etcdutil v1.3.5
	github.com/cybozu-go/log v1.6.0
	github.com/cybozu-go/well v1.10.0
	github.com/google/subcommands v1.2.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.4
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	sigs.k8s.io/yaml v1.2.0
)

go 1.13

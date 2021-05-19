module github.com/cybozu-go/etcdpasswd

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/cybozu-go/etcdutil v1.4.0
	github.com/cybozu-go/log v1.6.0
	github.com/cybozu-go/well v1.10.0
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.12.0
	github.com/spf13/cobra v1.1.3
	go.etcd.io/etcd v0.5.0-alpha.5.0.20210512015243-d19fbe541bf9
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	sigs.k8s.io/yaml v1.2.0
)

# Makefile for integration testing using ginkgo

# configuration variables
BRIDGE_ADDRESS = 10.0.0.1
HOST1 = 10.0.0.101
HOST2 = 10.0.0.102
HOST3 = 10.0.0.103
PLACEMAT = $(GOPATH)/bin/placemat
ETCD_VERSION = 3.3.8

# non-configuration variables
ETCD_ARCHIVE := etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz
SSH_PRIVKEY := ./itest_key
OUTPUT := ./output

GENERATED_FILES = $(OUTPUT)/etcd  $(OUTPUT)/etcdctl \
	$(OUTPUT)/etcdpasswd $(OUTPUT)/ep-agent \
	$(OUTPUT)/cluster.yml \
	$(OUTPUT)/user-data.yml \
	$(OUTPUT)/network1.yml \
	$(OUTPUT)/network2.yml \
	$(OUTPUT)/network3.yml

all: test

$(ETCD_ARCHIVE):
	curl -sSLf https://github.com/coreos/etcd/releases/download/v$(ETCD_VERSION)/etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz

$(OUTPUT)/etcd $(OUTPUT)/etcdctl: $(ETCD_ARCHIVE)
	mkdir -p $(OUTPUT)
	tar -C $(OUTPUT) -x -z --wildcards --strip-components 1 -f $< '*/$(notdir $@)'
	touch $@

$(OUTPUT)/etcdpasswd $(OUTPUT)/ep-agent:
	mkdir -p $(OUTPUT)
	cd ..; GOBIN=$(realpath $(OUTPUT)) go install ./cmd/$(notdir $@)

$(OUTPUT)/cluster.yml: cluster.yml
	mkdir -p $(OUTPUT)
	sed 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' $< > $@

$(OUTPUT)/user-data.yml: user-data.yml
	mkdir -p $(OUTPUT)
	cp $< $@

$(OUTPUT)/network1.yml: network.yml
	mkdir -p $(OUTPUT)
	sed -e 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' \
            -e 's/@HOST1@/$(HOST1)/' $< > $@

$(OUTPUT)/network2.yml: network.yml
	mkdir -p $(OUTPUT)
	sed -e 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' \
            -e 's/@HOST2@/$(HOST2)/' $< > $@

$(OUTPUT)/network3.yml: network.yml
	mkdir -p $(OUTPUT)
	sed -e 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' \
            -e 's/@HOST3@/$(HOST3)/' $< > $@

test: $(GENERATED_FILES)
	ginkgo

clean:
	rm -rf $(OUTPUT)

setup:
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega/...
	go get -u github.com/cybozu-go/placemat

.PHONY:	all clean setup test

# Makefile for integration testing using ginkgo

# configuration variables
BRIDGE_ADDRESS = 10.0.0.1
HOST1 = 10.0.0.101
HOST2 = 10.0.0.102
HOST3 = 10.0.0.103
PLACEMAT = $(GOPATH)/bin/placemat
GINKGO = $(GOPATH)/bin/ginkgo
ETCD_VERSION = 3.3.8
PLACEMAT_DATADIR = /var/scratch/placemat

export BRIDGE_ADDRESS HOST1 HOST2 HOST3 PLACEMAT GINKGO

# non-configuration variables
ETCD_ARCHIVE := etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz
SSH_PRIVKEY := ./itest_key
OUTPUT := ./output
UBUNTU_IMAGE := ubuntu-18.04-server-cloudimg-amd64.img

export SSH_PRIVKEY

GENERATED_FILES = $(OUTPUT)/etcd  $(OUTPUT)/etcdctl \
	$(OUTPUT)/etcdpasswd $(OUTPUT)/ep-agent \
	$(OUTPUT)/cluster.yml \
	$(OUTPUT)/user-data.yml \
	$(OUTPUT)/network1.yml \
	$(OUTPUT)/network2.yml \
	$(OUTPUT)/network3.yml \
	$(UBUNTU_IMAGE)

all: test

$(ETCD_ARCHIVE):
	curl -sSLf -o $@ https://github.com/coreos/etcd/releases/download/v$(ETCD_VERSION)/etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz

$(UBUNTU_IMAGE):
	curl -sSLf -o $@ https://cloud-images.ubuntu.com/releases/18.04/release/$@

$(OUTPUT)/etcd $(OUTPUT)/etcdctl: $(ETCD_ARCHIVE)
	mkdir -p $(OUTPUT)
	tar -C $(OUTPUT) -x -z --wildcards --strip-components 1 -f $< '*/$(notdir $@)'
	touch $@

$(OUTPUT)/etcdpasswd $(OUTPUT)/ep-agent:
	mkdir -p $(OUTPUT)
	cd ..; make all GOBIN=$(realpath $(OUTPUT))

$(OUTPUT)/cluster.yml: cluster.yml
	mkdir -p $(OUTPUT)
	sed 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' $< > $@

$(OUTPUT)/user-data.yml: user-data.yml
	mkdir -p $(OUTPUT)
	sed 's/@HOST1@/$(HOST1)/' $< > $@

$(OUTPUT)/network1.yml: network.yml
	mkdir -p $(OUTPUT)
	sed -e 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' \
            -e 's/@HOST@/$(HOST1)/' $< > $@

$(OUTPUT)/network2.yml: network.yml
	mkdir -p $(OUTPUT)
	sed -e 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' \
            -e 's/@HOST@/$(HOST2)/' $< > $@

$(OUTPUT)/network3.yml: network.yml
	mkdir -p $(OUTPUT)
	sed -e 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' \
            -e 's/@HOST@/$(HOST3)/' $< > $@

test: $(GENERATED_FILES)
	sudo rm -rf $(PLACEMAT_DATADIR)/volumes/host1
	sudo rm -rf $(PLACEMAT_DATADIR)/volumes/host2
	sudo rm -rf $(PLACEMAT_DATADIR)/volumes/host3
	./test.sh

placemat: $(GENERATED_FILES)
	sudo rm -rf $(PLACEMAT_DATADIR)/volumes/host1
	sudo rm -rf $(PLACEMAT_DATADIR)/volumes/host2
	sudo rm -rf $(PLACEMAT_DATADIR)/volumes/host3
	sudo -b $(PLACEMAT) output/cluster.yml
	@echo 'Placemat is running.  To kill it, do "sudo pkill placemat".'

test-light:
	$(GINKGO)

clean:
	rm -rf $(OUTPUT)

setup:
	GO111MODULE=off go get -u \
		github.com/onsi/ginkgo/ginkgo \
		github.com/cybozu-go/placemat/pkg/placemat

.PHONY:	all clean setup test

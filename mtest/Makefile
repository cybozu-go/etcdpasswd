# Makefile for integration testing using ginkgo

# configuration variables
BRIDGE_ADDRESS = 10.0.0.1
HOST1 = 10.0.0.101
HOST2 = 10.0.0.102
HOST3 = 10.0.0.103
PLACEMAT = /usr/bin/placemat
GINKGO = $(GOPATH)/bin/ginkgo
ETCD_VERSION = 3.3.12
PLACEMAT_DATADIR = /var/scratch/placemat
SUDO = sudo
SUITE = ./functions
ifeq ($(findstring /,$(SUITE)),)
  SUITE_PACKAGE = ./$(SUITE)
else
  SUITE_PACKAGE = $(SUITE)
endif

export BRIDGE_ADDRESS HOST1 HOST2 HOST3 PLACEMAT GINKGO SUITE_PACKAGE

# non-configuration variables
ETCD_ARCHIVE := etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz
SSH_PRIVKEY = $(realpath ./mtest_key)
OUTPUT := ./output
UBUNTU_IMAGE := ubuntu-18.04-server-cloudimg-amd64.img
ETCD := $(abspath $(OUTPUT))/etcd
ETCDCTL := $(abspath $(OUTPUT))/etcdctl
ETCDPASSWD := $(abspath $(OUTPUT))/etcdpasswd
EPAGENT := $(abspath $(OUTPUT))/ep-agent

export SSH_PRIVKEY ETCD ETCDCTL ETCDPASSWD EPAGENT

TEST_DEPS = \
	$(OUTPUT)/etcd \
	$(OUTPUT)/etcdctl \
	$(OUTPUT)/etcdpasswd \
	$(OUTPUT)/ep-agent

PLACEMAT_DEPS = \
	$(OUTPUT)/cluster.yml \
	$(OUTPUT)/user-data.yml \
	$(OUTPUT)/network1.yml \
	$(OUTPUT)/network2.yml \
	$(OUTPUT)/network3.yml \
	$(UBUNTU_IMAGE)

all:
	@echo "Specify one of these targets:"
	@echo
	@echo "    placemat - run VMs in placemat."
	@echo "    stop     - stop placemat."
	@echo "    test     - run test in placemat env."
	@echo "    setup    - install dependencies."
	@echo
	@echo "'make test' accepts SUITE."

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
	cd ..; GOBIN=$(realpath $(OUTPUT)) go install ./pkg/$(notdir $@)

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

# Targets
test: $(TEST_DEPS)
	./test.sh "$(TARGET)"

placemat: $(PLACEMAT_DEPS) $(TEST_DEPS)
	sudo rm -rf $(PLACEMAT_DATADIR)
	sudo rm -rf $(VM_SHARED_DIR)
	sudo systemd-run --unit=placemat.service $(PLACEMAT) -enable-virtfs $(abspath output/cluster.yml)
	@echo 'Placemat is running.  To kill it, do "make stop".'

stop:
	sudo systemctl stop placemat.service
	sudo systemctl disable placemat.service
	sudo rm /run/systemd/transient/placemat.service
	sudo systemctl daemon-reload
	sudo systemctl reset-failed placemat.service

clean:
	rm -rf $(OUTPUT)

setup:
	go install github.com/onsi/ginkgo/ginkgo

.PHONY: all test placemat stop clean setup

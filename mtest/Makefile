# Makefile for integration testing using ginkgo

# configuration variables
BRIDGE_ADDRESS = 10.0.0.1
HOST1 = 10.0.0.101
HOST2 = 10.0.0.102
HOST3 = 10.0.0.103
PLACEMAT = /usr/bin/placemat2
ETCD_VERSION = 3.5.7
PLACEMAT_DATADIR = /var/scratch/placemat
SUDO = sudo

export HOST1 HOST2 HOST3 PLACEMAT

# non-configuration variables
ETCD_ARCHIVE := etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz
SSH_PRIVKEY = $(realpath ./mtest_key)
OUTPUT := ./output
DATA_DIR := $(abspath $(OUTPUT))/data
UBUNTU_IMAGE := ubuntu-22.04-server-cloudimg-amd64.img
ETCD := $(DATA_DIR)/etcd
ETCDCTL := $(DATA_DIR)/etcdctl
ETCDPASSWD := $(DATA_DIR)/etcdpasswd
EPAGENT := $(DATA_DIR)/ep-agent

export SSH_PRIVKEY ETCD ETCDCTL ETCDPASSWD EPAGENT

TEST_DEPS = \
	$(DATA_DIR)/etcd \
	$(DATA_DIR)/etcdctl \
	$(DATA_DIR)/etcdpasswd \
	$(DATA_DIR)/ep-agent

PLACEMAT_DEPS = \
	$(OUTPUT)/cluster.yml \
	$(OUTPUT)/user-data.yml \
	$(OUTPUT)/network1.yml \
	$(OUTPUT)/network2.yml \
	$(OUTPUT)/network3.yml \
	$(UBUNTU_IMAGE)

.PHONY: all
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
	curl -sSLf -o $@ https://github.com/etcd-io/etcd/releases/download/v$(ETCD_VERSION)/etcd-v$(ETCD_VERSION)-linux-amd64.tar.gz

$(UBUNTU_IMAGE):
	curl -sSLf -o $@ https://cloud-images.ubuntu.com/releases/22.04/release/$@

$(DATA_DIR)/etcd $(DATA_DIR)/etcdctl: $(ETCD_ARCHIVE)
	mkdir -p $(DATA_DIR)
	tar -C $(DATA_DIR) -x -z --wildcards --strip-components 1 -f $< '*/$(notdir $@)'
	touch $@

$(DATA_DIR)/etcdpasswd $(DATA_DIR)/ep-agent:
	mkdir -p $(DATA_DIR)
	cd ..; GOBIN=$(realpath $(DATA_DIR)) go install ./pkg/$(notdir $@)

$(OUTPUT)/cluster.yml: cluster.yml
	mkdir -p $(OUTPUT)
	sed -e 's/@BRIDGE_ADDRESS@/$(BRIDGE_ADDRESS)/' \
	    -e 's#@DATA_DIR@#$(DATA_DIR)#' $< > $@

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
.PHONY: test
test: $(TEST_DEPS)
	./test.sh "$(TARGET)"

.PHONY: placemat
placemat: $(PLACEMAT_DEPS) $(TEST_DEPS)
	sudo rm -rf $(PLACEMAT_DATADIR)
	sudo rm -rf $(VM_SHARED_DIR)
	sudo systemd-run --unit=placemat.service $(PLACEMAT) $(abspath output/cluster.yml)
	@echo 'Placemat is running.  To kill it, do "make stop".'

.PHONY: stop
stop:
	sudo systemctl stop placemat.service
	sudo systemctl disable placemat.service
	sudo rm /run/systemd/transient/placemat.service
	sudo systemctl daemon-reload
	sudo systemctl reset-failed placemat.service

.PHONY: clean
clean:
	rm -rf $(OUTPUT)

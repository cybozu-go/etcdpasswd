# Makefile for etcdpasswd

PACKAGES := fakeroot
DOC_DIR := debian/usr/share/doc/etcdpasswd
CONTROL := debian/DEBIAN/control
SUDO = sudo

ETCD_VER=3.6.11
ETCD_SHA256=8756f7a4eaf921668a83de0bf13c0f65cae9186a165696e3ae8396afe6f557ed

# Test tools
BIN_DIR := $(shell pwd)/bin
ETCD := $(BIN_DIR)/etcd

GOLANGCI_LINT := go tool -modfile=tools/go.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint

all: test

.PHONY: check-generate
check-generate:
	go mod tidy
	git diff --exit-code --name-only

.PHONY: test
test: lint
	go test -race -count=1 -v ./...

.PHONY: lint
lint:
	$(GOLANGCI_LINT) run -vvv

.PHONY: lint-fix
lint-fix:
	$(GOLANGCI_LINT) run -vvv --fix

$(CONTROL): control
	sed 's/@VERSION@/$(patsubst v%,%,$(VERSION))/' $< > $@

.PHONY: deb
deb: $(CONTROL)
	mkdir -p debian/usr/bin
	GOBIN=$(CURDIR)/debian/usr/bin CGO_ENABLED=0 go install -ldflags="-s -w" ./pkg/etcdpasswd
	mkdir -p debian/usr/sbin
	GOBIN=$(CURDIR)/debian/usr/sbin CGO_ENABLED=0 go install -ldflags="-s -w" ./pkg/ep-agent
	mkdir -p $(DOC_DIR)
	cp config.yml.example README.md LICENSE $(DOC_DIR)
	cp pkg/etcdpasswd/USAGE.md $(DOC_DIR)/etcdpasswd.md
	mkdir -p debian/lib/systemd/system
	cp pkg/ep-agent/ep-agent.service debian/lib/systemd/system
	chmod -R g-w debian
	fakeroot dpkg-deb --build debian .

.PHONY: test-tools
test-tools: $(STATICCHECK) $(CUSTOM_CHECKER) $(ETCD)

.PHONY: clean
clean:
	rm -f *.deb
	rm -rf $(CONTROL) debian/usr debian/lib
	rm -rf $(BIN_DIR)

.PHONY: setup
setup:
	$(SUDO) apt-get update
	$(SUDO) apt-get -y --no-install-recommends install $(PACKAGES)

$(ETCD):
	mkdir -p $(BIN_DIR)
	curl -sL https://github.com/etcd-io/etcd/releases/download/v${ETCD_VER}/etcd-v${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz
	echo "${ETCD_SHA256}  /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz" | sha256sum -c -
	mkdir /tmp/etcd
	tar xzvf /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd --strip-components=1
	$(SUDO) mv /tmp/etcd/etcd $(ETCD)
	rm -rf /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz /tmp/etcd

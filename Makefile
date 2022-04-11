# Makefile for etcdpasswd

PACKAGES := fakeroot
DOC_DIR := debian/usr/share/doc/etcdpasswd
CONTROL := debian/DEBIAN/control
SUDO = sudo

ETCD_VER=3.5.1

# Test tools
BIN_DIR := $(shell pwd)/bin
STATICCHECK := $(BIN_DIR)/staticcheck
NILERR := $(BIN_DIR)/nilerr
CUSTOM_CHECKER := $(BIN_DIR)/custom-checker
ETCD := $(BIN_DIR)/etcd

all: test

.PHONY: check-generate
check-generate:
	go mod tidy
	git diff --exit-code --name-only

.PHONY: test
test:
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"
	$(STATICCHECK) ./...
	$(NILERR) ./...
	test -z "$$($(CUSTOM_CHECKER) -restrictpkg.packages=html/template,log $$(go list ./...) 2>&1 | tee /dev/stderr)"
	go test -race -count=1 -v ./...
	go vet ./...

$(CONTROL): control
	sed 's/@VERSION@/$(patsubst v%,%,$(VERSION))/' $< > $@

.PHONY: deb
deb: $(CONTROL)
	mkdir -p debian/usr/bin
	GOBIN=$(CURDIR)/debian/usr/bin go install ./pkg/etcdpasswd
	mkdir -p debian/usr/sbin
	GOBIN=$(CURDIR)/debian/usr/sbin go install ./pkg/ep-agent
	mkdir -p debian/usr/share/doc/etcdpasswd
	cp config.yml.example README.md LICENSE $(DOC_DIR)
	cp pkg/etcdpasswd/USAGE.md $(DOC_DIR)/etcdpasswd.md
	mkdir -p debian/lib/systemd/system
	cp pkg/ep-agent/ep-agent.service debian/lib/systemd/system
	chmod -R g-w debian
	fakeroot dpkg-deb --build debian .

.PHONY: test-tools
test-tools: $(STATICCHECK) $(NILERR) $(CUSTOM_CHECKER) $(ETCD)

.PHONY: clean
clean:
	rm -f *.deb
	rm -rf $(CONTROL) debian/usr debian/lib
	rm -rf $(BIN_DIR)

.PHONY: setup
setup:
	$(SUDO) apt-get update
	$(SUDO) apt-get -y --no-install-recommends install $(PACKAGES)

$(STATICCHECK):
	mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install honnef.co/go/tools/cmd/staticcheck@latest

$(NILERR):
	mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/gostaticanalysis/nilerr/cmd/nilerr@latest

$(CUSTOM_CHECKER):
	mkdir -p $(BIN_DIR)
	GOBIN=$(BIN_DIR) go install github.com/cybozu/neco-containers/golang/analyzer/cmd/custom-checker@latest

$(ETCD):
	mkdir -p $(BIN_DIR)
	curl -sL https://github.com/etcd-io/etcd/releases/download/v${ETCD_VER}/etcd-v${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz
	mkdir /tmp/etcd
	tar xzvf /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd --strip-components=1
	$(SUDO) mv /tmp/etcd/etcd $(ETCD)
	rm -rf /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz /tmp/etcd

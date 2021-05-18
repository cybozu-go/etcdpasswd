# Makefile for etcdpasswd

PACKAGES := fakeroot
DOCDIR := debian/usr/share/doc/etcdpasswd
CONTROL := debian/DEBIAN/control
SUDO = sudo

ETCD_VER=3.4.16

all: test

.PHONY: test
test: test-tools
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"
	staticcheck ./...
	test -z "$$(nilerr ./... 2>&1 | tee /dev/stderr)"
	test -z "$$(custom-checker -restrictpkg.packages=html/template,log $$(go list ./...) 2>&1 | tee /dev/stderr)"
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
	cp config.yml.example README.md LICENSE $(DOCDIR)
	cp pkg/etcdpasswd/USAGE.md $(DOCDIR)/etcdpasswd.md
	mkdir -p debian/lib/systemd/system
	cp pkg/ep-agent/ep-agent.service debian/lib/systemd/system
	chmod -R g-w debian
	fakeroot dpkg-deb --build debian .

.PHONY: test-tools
test-tools: staticcheck nilerr goimports custom-checker etcd

.PHONY: clean
clean:
	rm -f *.deb
	rm -rf $(CONTROL) debian/usr debian/lib

.PHONY: setup
setup:
	$(SUDO) apt-get update
	$(SUDO) apt-get -y --no-install-recommends install $(PACKAGES)

.PHONY: staticcheck
staticcheck:
	if ! which staticcheck >/dev/null; then \
		env GOFLAGS= go install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi

.PHONY: nilerr
nilerr:
	if ! which nilerr >/dev/null; then \
		env GOFLAGS= go install github.com/gostaticanalysis/nilerr/cmd/nilerr@latest; \
	fi

.PHONY: goimports
goimports:
	if ! which goimports >/dev/null; then \
		env GOFLAGS= go install golang.org/x/tools/cmd/goimports@latest; \
	fi

.PHONY: custom-checker
custom-checker:
	if ! which custom-checker >/dev/null; then \
		env GOFLAGS= go install github.com/cybozu/neco-containers/golang/analyzer/cmd/custom-checker@latest; \
	fi

.PHONY: etcd
etcd:
	if ! which etcd >/dev/null; then \
		curl -L https://github.com/etcd-io/etcd/releases/download/v${ETCD_VER}/etcd-v${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz; \
		mkdir /tmp/etcd; \
		tar xzvf /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd --strip-components=1; \
		$(SUDO) mv /tmp/etcd/etcd /usr/local/bin/; \
		rm -rf /tmp/etcd-v${ETCD_VER}-linux-amd64.tar.gz /tmp/etcd; \
	fi

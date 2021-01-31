# Makefile for etcdpasswd

PACKAGES := fakeroot
DOCDIR := debian/usr/share/doc/etcdpasswd
CONTROL := debian/DEBIAN/control
SUDO = sudo

all: test

test:
	test -z "$$(gofmt -s -l . | tee /dev/stderr)"
	staticcheck ./...
	test -z "$$(nilerr ./... 2>&1 | tee /dev/stderr)"
	test -z "$$(custom-checker -restrictpkg.packages=html/template,log $$(go list ./...) 2>&1 | tee /dev/stderr)"
	go test -race -count=1 -v ./...
	go vet ./...

$(CONTROL): control
	sed 's/@VERSION@/$(patsubst v%,%,$(VERSION))/' $< > $@

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

clean:
	rm -f *.deb
	rm -rf $(CONTROL) debian/usr debian/lib

setup:
	$(SUDO) apt-get update
	$(SUDO) apt-get -y --no-install-recommends install $(PACKAGES)

.PHONY: all test deb clean setup

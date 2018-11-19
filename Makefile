# Makefile for etcdpasswd

PACKAGES := fakeroot
DOCDIR := debian/usr/share/doc/etcdpasswd
CONTROL := debian/DEBIAN/control

# for Go
GOFLAGS = -mod=vendor
export GOFLAGS

all:
	go install ./...

test:
	test -z "$$(gofmt -s -l . | grep -v '^vendor' | tee /dev/stderr)"
	golint -set_exit_status $$(go list -mod=vendor ./... | grep -v /vendor/)
	go test -race -count=1 -v ./...
	go vet ./...

$(CONTROL): control
	sed 's/@VERSION@/$(patsubst v%,%,$(VERSION))/' $< > $@

deb: $(CONTROL)
	mkdir -p debian/usr/bin
	GOBIN=$(CURDIR)/debian/usr/bin go install ./cmd/etcdpasswd
	mkdir -p debian/usr/sbin
	GOBIN=$(CURDIR)/debian/usr/sbin go install ./cmd/ep-agent
	mkdir -p debian/usr/share/doc/etcdpasswd
	cp etcdpasswd.yml.example README.md LICENSE $(DOCDIR)
	cp cmd/etcdpasswd/USAGE.md $(DOCDIR)/etcdpasswd.md
	mkdir -p debian/lib/systemd/system
	cp cmd/ep-agent/ep-agent.service debian/lib/systemd/system
	chmod -R g-w debian
	fakeroot dpkg-deb --build debian .

clean:
	rm -f *.deb
	rm -rf $(CONTROL) debian/usr debian/lib

setup:
	sudo apt-get -y --no-install-recommends install $(PACKAGES)

.PHONY: all test deb clean setup

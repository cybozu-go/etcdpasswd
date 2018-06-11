PACKAGES := fakeroot
DOCDIR := debian/usr/share/doc/etcdpasswd

all:
	go get -d ./...
	go install ./...

test:
	go get -d -t ./...
	go test -race -count 1 -v ./...
	test -z "$$(goimports -d . | tee /dev/stderr)"
	golint -set_exit_status ./...

deb:
	go get -d ./...
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
	rm -rf debian/usr debian/lib

setup:
	sudo apt-get -y --no-install-recommends install $(PACKAGES)

.PHONY: all test deb clean setup

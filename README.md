[![GitHub release](https://img.shields.io/github/release/cybozu-go/etcdpasswd.svg?maxAge=60)][releases]
[![GoDoc](https://godoc.org/github.com/cybozu-go/etcdpasswd?status.svg)][godoc]
[![CI](https://github.com/cybozu-go/etcdpasswd/workflows/main/badge.svg)](https://github.com/cybozu-go/etcdpasswd/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/cybozu-go/etcdpasswd)](https://goreportcard.com/report/github.com/cybozu-go/etcdpasswd)

etcdpasswd
==========

etcdpasswd manages Linux users and groups with a central database on etcd.
This repository provides following two programs:

* `ep-agent`: a background service that watches etcd database and synchronize Linux users/groups.
* `etcdpasswd`: CLI tool to edit the central database on etcd.

Build
-----

```console
$ go install github.com/cybozu-go/etcdpasswd/...
```

Installation
------------

1. Prepare an etcd cluster.

2. Create `/etc/etcdpasswd/config.yml`.

    This file provides parameters to connect to the etcd cluster.
    A sample configuration looks like this:

    ```yaml
    endpoints:
      - http://12.34.56.78:2379
    username: cybozu
    password: xxxxxxxx

    tls-cert-file: /etc/etcdpasswd/etcd.crt
    tls-key-file: /etc/etcdpasswd/etcd.key
    ```

3. Run `ep-agent`.

    A sample systemd unit file is available at [cmd/ep-agent/ep-agent.service](cmd/ep-agent/ep-agent.service).
    Use it to run `ep-agent` as a systemd service as follows:

    ```console
    $ sudo cp $GOPATH/bin/ep-agent /usr/local/sbin
    $ sudo cp ep-agent.service /etc/systemd/system
    $ sudo systemctl daemon-reload
    $ sudo systemctl enable ep-agent.service
    $ sudo systemctl start ep-agent.service
    ```

4. Use `etcdpasswd` to initialize the database.

    ```console
    $ etcdpasswd set start-uid 2000
    $ etcdpasswd set start-gid 2000
    $ etcdpasswd set default-group cybozu
    $ etcdpasswd set default-groups sudo,adm
    ```

Usage
-----

See [pkg/etcdpasswd/USAGE.md](pkg/etcdpasswd/USAGE.md).

Specifications
--------------

Read [docs/spec.md](docs/spec.md).

License
-------

MIT

[releases]: https://github.com/cybozu-go/etcdpasswd/releases
[godoc]: https://godoc.org/github.com/cybozu-go/etcdpasswd

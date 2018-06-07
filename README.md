etcdpasswd
==========

etcdpasswd manages Linux users and groups with a central database on etcd.
This repository provides following two programs:

* `ep-agent`: a background service that watches etcd database and synchronize Linux users/groups.
* `etcdpasswd`: CLI tool to edit the central database on etcd.

Build
-----

```console
$ go get -u github.com/cybozu-go/etcdpasswd/...
```

Usage
-----

TODO

Specifications
--------------

Read [SPEC.md](SPEC.md).

License
-------

MIT

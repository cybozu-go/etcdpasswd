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

1. Create config file as `/etc/etcdpasswd.yml`.

    ```yaml
    servers:
      - http://12.34.56.78:2379
    username: cybozu
    password: xxxxxxxx
    ```

Specifications
--------------

Read [SPEC.md](SPEC.md).

License
-------

MIT

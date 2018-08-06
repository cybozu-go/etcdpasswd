Specifications
==============

`/etc/etcdpasswd.yml`
---------------------

This file provides connecting to the etcd cluster.
Parameters are defined by [cybozu-go/etcdutil](https://github.com/cybozu-go/etcdutil), and not shown below will use default values of the etcdutil.

Name     | Type   | Required | Description
-------- | ------ | -------- | -----------
`prefix` | string | No       | Key prefix of etcd objects.  Default is `/passwd/`.

etcd schema
-----------

etcdpasswd stores following data in etcd.

### Prefix

The default prefix of keys in etcd is `/passwd/`.

### Configuration

`<prefix>/config` holds etcdpasswd configurations in JSON format like this:

```json
{
  "start-uid": 2000,
  "start-gid": 2000,
  "default-group": "cybozu",
  "default-groups": ["sudo", "adm"],
  "default-shell": "/bin/bash"
}
```

Key            | Type            | Description
---            | ----            | -----------
start-uid      | int             | Starting UID for users managed by etcdpasswd.
start-gid      | int             | Starting GID for groups managed by etcdpasswd.
default-group  | string          | Default primary group for new users.
default-groups | array of string | Default supplementary groups for new users.
default-shell  | string          | Default shell program for new users.

`start-uid` and `start-gid` must be set before adding new users or groups.

### Last UID

`<prefix>/last-uid` holds the last used UID.
If this key exists, the next user will have `last-uid + 1` as their UID.
If this key does not exist, the next user will have `start-uid` as their UID.

```json
2003
```

### Last GID

`<prefix>/last-gid` holds the last used GID.

```json
2004
```

### User information

`<prefix>/users/<name>` holds user information.

```json
{
  "name": "cybozu",
  "uid": 2006,
  "display-name": "Cy Bozu",
  "group": "cybozu",
  "groups": ["sudo", "adm"],
  "shell": "/bin/bash",
  "public-keys": ["public-key-1", "public-key-2"]
}
```

### Group information

`<prefix>/groups/<name>` holds the GID of `<name>` group.

```json
2008
```

### Deleted users

`<prefix>/deleted-users/<name>` indicates that `<name>` user has been deleted.

When a new user with the same name is to be created, this key will be removed in the same transaction.

### Deleted groups

`<prefix>/deleted-groups/<name>` indicates that `<name>` group has been deleted.

When a new group with the same name is to be created, this key will be removed in the same transaction.

### Unmanaged users to be locked

If `<prefix>/locked/<name>` key exists, `<name>` user will automatically be locked by ep-agent.

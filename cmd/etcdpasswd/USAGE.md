etcdpasswd usage
================

## Synopsis

`etcdpasswd [OPTIONS] SUBCOMMAND`

Options:
- `-config`: configuration file path (default: `/etc/etcdpasswd.yml`)

Subcommands:
- `set CONFIG VALUE`
- `get CONFIG`
- `lock NAME`
- `unlock NAME`
- `user list`
- `user get NAME`
- `user add [OPTIONS] NAME`
- `user update [OPTIONS] NAME`
- `user remove NAME`
- `cert list NAME`
- `cert add NAME [FILE]`
- `cert remove NAME INDEX`
- `group list`
- `group add NAME`
- `group remove NAME`

## `set CONFIG VALUE` and `get CONFIG`

Set/get etcdpasswd configurations.

Configurations are:
- `start-uid`: starting UID for users managed by etcdpasswd.
- `start-gid`: starting GID for groups managed by etcdpasswd.
- `default-group`: default primary group name.
- `default-groups`: comma-separated list of default supplementary group names.
- `default-shell`: default shell program.


## `lock NAME` and `unlock NAME`

`lock` adds a user to etcd database to instruct ep-agent to lock the user by `passwd -l`.

`unlock` removes the user from etcd database.  It does not unlock the user's password.
Administrators need to manually unlock the user by `passwd -u`.

## `user`

## `cert`

## `group`

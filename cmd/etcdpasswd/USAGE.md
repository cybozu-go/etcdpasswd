etcdpasswd usage
================

## Synopsis

`etcdpasswd [OPTIONS] SUBCOMMAND`

Options:
- `-config`: configuration file path (default: `/etc/etcdpasswd.yml`)

Subcommands:
- `set CONFIG VALUE`
- `get CONFIG`
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
- `locker list`
- `locker add NAME`
- `locker remove NAME`

## `set CONFIG VALUE` and `get CONFIG`

Set/get etcdpasswd configurations.

Configurations are:
- `start-uid`: starting UID for users managed by etcdpasswd.
- `start-gid`: starting GID for groups managed by etcdpasswd.
- `default-group`: default primary group name.
- `default-groups`: comma-separated list of default supplementary group names.
- `default-shell`: default shell program.

## `user`

## `cert`

## `group`

## `locker`

Lock passwords to prevent logins using password.

`locker add NAME` adds the named user to the list of password-locked users in the database.
The list is watched by `ep-agent` who will lock passwords of users in the list by `passwd -l`.

The user is not limited to users added by `etcdpasswd`; local users including `root` may
be locked by this.

`locker remove NAME` removes the name from the list, but `ep-agent` does *not* unlock
passwords.  Unlocking passwords should be done by administrators later on.

`locker list` shows names in the list.

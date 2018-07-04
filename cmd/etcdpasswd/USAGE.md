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

Manage users.

### `user add`

`user add [OPTIONS] NAME` adds a new managed user to the database.
`ep-agent` watches the database to actually add the user to each Unix system.

Available options are:
- `display`: the user's display name.  The default is "".
- `group`: primary group name of the user.  The default is taken from the configuration.
- `groups`: comma-separated list of the user's supplementary group names.  The default is taken from the configuration, or none if not configured.
- `shell`: the user's login shell.  The default is taken from the configuration, or `/bin/bash` if not configured.

The UID of the new user is taken in sequence from `start-uid` in the configuration.

`user add` returns an error if, but not only if:
- `NAME` is used by another managed user.
- `NAME` is not a valid user name.
- at least one of the specified group names is not a valid group name.
- primary group name is not given either by the command line or the configuration.
- `start-uid` has not been set.

`user add` does NOT return an error in the following cases, but `ep-agent` will fail to reflect the change:
- UID is used by an unmanaged user.
- at least one of the specified groups does not exist on the system.

It is NOT an error for `user add` or `ep-agent` that `NAME` is used by an unmanaged user.
`ep-agent` first removes the unmanaged user and then creates a new managed user on each system.

### `user update`

`user update [OPTION] NAME` updates an existing managed user in the database.
`ep-agent` watches the database to actually update the user on each Unix system.

The same options are available as `user add`.
Only the attributes specified by options with non-empty values are updated.

`user update` returns an error if, but not only if:
- `NAME` is not a managed user's name.
- at least one of the specified group names is not a valid group name.
- the user is being updated concurrently.

### `user remove`

`user remove NAME` removes an existing managed user from the database.
`ep-agent` watches the database to actually remove the user from each Unix system.

The home directory of the user is also removed.

`user remove` returns an error if `NAME` is not a managed user's name.

### `user list`

`user list` lists all user names registered in the database.
The result is sorted alphabetically.

### `user get`

`user get NAME` shows the attributes of the named user registered in the database.

`user get` returns an error if `NAME` is not a managed user's name.

## `cert`

Manage SSH public keys of users.

`cert add NAME [FILE]` adds an SSH public key for the named user into the database.
`ep-agent` watches the database to actually add the key to the user's `.ssh/authorized_keys` on each Unix system.
The SSH public key is read from `FILE`, or from stdin if `FILE` is not specified.

`cert remove NAME INDEX` removes an SSH public key for the named user from the database.
`ep-agent` watches the database to actually remove the key from the user's `.ssh/authorized_keys` on each Unix system.
The key to be removed is specified by `INDEX`, which is shown by `cert list`.

`cert list NAME` shows a summary of the SSH public keys for the named user in the database.
The summary consists of lines of public key information, one line per public key,
and each line is formatted as `<index>: <comment> (<type>)`, e.g. `0: foo@localhost (ssh-rsa)`.

## `group`

Manage groups.

### `group add`

`group add NAME` adds a new managed group to the database.
`ep-agent` watches the database to actually add the group to each Unix system.

The GID of the new group is taken in sequence from `start-gid` in the configuration.

`group add` returns an error if, but not only if:
- `NAME` is used by another managed group.
- `NAME` is not a valid group name.
- `start-gid` has not been set.

`group add` does NOT return an error in the following case, but `ep-agent` will fail to reflect the change:
- GID is used by an unmanaged group.

It is NOT an error for `group add` or `ep-agent` that `NAME` is used by an unmanaged group.
`ep-agent` first removes the unmanaged group and then creates a new managed group on each system.

### `group remove`

`group remove NAME` removes an existing managed group from the database.
`ep-agent` watches the database to actually remove the group from each Unix system.

`group remove` returns an error if `NAME` is not a managed group's name.

`group remove` does NOT return an error in the following case, but `ep-agent` will fail to reflect the change:
- the group is referred as a primary group by some user.

### `group list`

`group list` lists all group names registered in the database.
The result is sorted alphabetically.

## `locker`

Lock passwords to prevent logins using password.

`locker add NAME` adds the named user to the list of password-locked users in the database.
The list is watched by `ep-agent` who will lock passwords of users in the list by `passwd -l`.

The user is not limited to users added by `etcdpasswd`; local users including `root` may
be locked by this.

`locker remove NAME` removes the name from the list, but `ep-agent` does *not* unlock
passwords.  Unlocking passwords should be done by administrators later on.

`locker list` shows names in the list.

## Valid name

A user/group name is valid if it matches the pattern of `^[a-z][-a-z0-9_]*$` and it does not conflict with system names such as "root" or "nobody".

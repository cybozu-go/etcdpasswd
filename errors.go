package etcdpasswd

type errString string

func (e errString) Error() string {
	return string(e)
}

const (
	// ErrCASFailure indicates compare-and-swap failure.
	ErrCASFailure = errString("conflicted")

	// ErrNotFound indicates an object was not found in the database.
	ErrNotFound = errString("not found")

	// ErrExists indicates that an object with the same key already exists.
	ErrExists = errString("already exists")
)

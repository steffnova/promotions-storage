package storage

// Option assign new capabilities to existing storage
type Option func(Storage) Storage

// WithOptions returns new [Storage] with added option. The storage
// parameter specifies storage to which options will be added. The
// options parameter represent list of [Option] that will be assigned
// to storage.
func WithOptions(storage Storage, options ...Option) Storage {
	for _, option := range options {
		storage = option(storage)
	}

	return storage
}

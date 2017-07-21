package gorocksdb

// Iterator provides a way to seek to specific keys and iterate through
// the keyspace from that point, as well as access the values of those keys.
//
type Iterator interface {
	// Valid returns false only when an Iterator has iterated past either the
	// first or the last key in the database.
	Valid() bool

	// ValidForPrefix returns false only when an Iterator has iterated past the
	// first or the last key in the database or the specified prefix.
	ValidForPrefix(prefix []byte) bool

	// Key returns the key the iterator currently holds.
	Key() *Slice

	// Value returns the value in the database the iterator currently holds.
	Value() *Slice

	// Next moves the iterator to the next sequential key in the database.
	Next()

	// Prev moves the iterator to the previous sequential key in the database.
	Prev()

	// SeekToFirst moves the iterator to the first key in the database.
	SeekToFirst()

	// SeekToLast moves the iterator to the last key in the database.
	SeekToLast()

	// Seek moves the iterator to the position greater than or equal to the key.
	Seek(key []byte)

	// SeekForPrev moves the iterator to the last key that less than or equal
	// to the target key, in contrast with Seek.
	SeekForPrev(key []byte)

	// Err returns nil if no errors happened during iteration, or the actual
	// error otherwise.
	Err() error

	// Close closes the iterator.
	Close()
}

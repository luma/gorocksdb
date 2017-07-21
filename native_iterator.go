package gorocksdb

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"
import (
	"bytes"
	"errors"
	"unsafe"
)

// NativeIterator provides a way to seek to specific keys and iterate through
// the keyspace from that point, as well as access the values of those keys.
//
// For example:
//
//      it := db.NewIterator(readOpts)
//      defer it.Close()
//
//      it.Seek([]byte("foo"))
//		for ; it.Valid(); it.Next() {
//          fmt.Printf("Key: %v Value: %v\n", it.Key().Data(), it.Value().Data())
// 		}
//
//      if err := it.Err(); err != nil {
//          return err
//      }
//
type NativeIterator struct {
	c *C.rocksdb_iterator_t
}

// NewNativeIterator creates a NativeIterator object.
func NewNativeIterator(c unsafe.Pointer) *NativeIterator {
	return &NativeIterator{(*C.rocksdb_iterator_t)(c)}
}

// Valid returns false only when an Iterator has iterated past either the
// first or the last key in the database.
func (iter *NativeIterator) Valid() bool {
	return C.rocksdb_iter_valid(iter.c) != 0
}

// ValidForPrefix returns false only when an Iterator has iterated past the
// first or the last key in the database or the specified prefix.
func (iter *NativeIterator) ValidForPrefix(prefix []byte) bool {
	if C.rocksdb_iter_valid(iter.c) == 0 {
		return false
	}

	key := iter.Key()
	result := bytes.HasPrefix(key.Data(), prefix)
	key.Free()
	return result
}

// Key returns the key the iterator currently holds.
func (iter *NativeIterator) Key() *Slice {
	var cLen C.size_t
	cKey := C.rocksdb_iter_key(iter.c, &cLen)
	if cKey == nil {
		return nil
	}
	return &Slice{cKey, cLen, true}
}

// Value returns the value in the database the iterator currently holds.
func (iter *NativeIterator) Value() *Slice {
	var cLen C.size_t
	cVal := C.rocksdb_iter_value(iter.c, &cLen)
	if cVal == nil {
		return nil
	}
	return &Slice{cVal, cLen, true}
}

// Next moves the iterator to the next sequential key in the database.
func (iter *NativeIterator) Next() {
	C.rocksdb_iter_next(iter.c)
}

// Prev moves the iterator to the previous sequential key in the database.
func (iter *NativeIterator) Prev() {
	C.rocksdb_iter_prev(iter.c)
}

// SeekToFirst moves the iterator to the first key in the database.
func (iter *NativeIterator) SeekToFirst() {
	C.rocksdb_iter_seek_to_first(iter.c)
}

// SeekToLast moves the iterator to the last key in the database.
func (iter *NativeIterator) SeekToLast() {
	C.rocksdb_iter_seek_to_last(iter.c)
}

// Seek moves the iterator to the position greater than or equal to the key.
func (iter *NativeIterator) Seek(key []byte) {
	cKey := byteToChar(key)
	C.rocksdb_iter_seek(iter.c, cKey, C.size_t(len(key)))
}

// SeekForPrev moves the iterator to the last key that less than or equal
// to the target key, in contrast with Seek.
func (iter *NativeIterator) SeekForPrev(key []byte) {
	cKey := byteToChar(key)
	C.rocksdb_iter_seek_for_prev(iter.c, cKey, C.size_t(len(key)))
}

// Err returns nil if no errors happened during iteration, or the actual
// error otherwise.
func (iter *NativeIterator) Err() error {
	var cErr *C.char
	C.rocksdb_iter_get_error(iter.c, &cErr)
	if cErr != nil {
		defer C.free(unsafe.Pointer(cErr))
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// Close closes the iterator.
func (iter *NativeIterator) Close() {
	C.rocksdb_iter_destroy(iter.c)
	iter.c = nil
}

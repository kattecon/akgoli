// Package utils provides small focused helpers for string handling, masking,
// cryptographic random IDs, constant-time comparison, buffer pooling,
// sentinel errors, rune-safe string slicing, ASCII validation, and in-place
// slice operations.
package utils

// ConstError is a string type that implements the error interface. Because
// its underlying type is string, a ConstError value can be declared as a
// package-level const and used as a sentinel error with == and errors.Is.
type ConstError string

func (e ConstError) Error() string { return string(e) }

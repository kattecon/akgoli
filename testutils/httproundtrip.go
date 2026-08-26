// Package testutils provides test helpers for buffered log capture, panic
// recovery, standard-stream redirection, and HTTP transport fakes.
package testutils

import "net/http"

// RoundTripFunc is a function type that implements http.RoundTripper. It lets
// tests supply an HTTP transport as a plain function without defining a struct.
type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

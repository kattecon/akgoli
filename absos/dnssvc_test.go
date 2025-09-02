package absos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnsSvc(t *testing.T) {
	svc := NewDnsSvc()

	assert.Equal(t, svc, NewDnsSvc())

	// Test with localhost.
	ips, err := svc.LookupIP("localhost")
	assert.Nil(t, err)
	assert.NotEmpty(t, ips)

	// Test with localhost.
	ips, err = svc.LookupIP("1.2.3.4")
	assert.Nil(t, err)
	assert.NotEmpty(t, ips)

	// Test with invalid hostname.
	_, err = svc.LookupIP("......invalid-hostname-that-should-not-exist.invalid.........")
	assert.NotNil(t, err)
}

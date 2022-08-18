package testutils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	xResp := &http.Response{}
	xReq := httptest.NewRequest(http.MethodGet, "/", nil)
	xErr := errors.New("xx")

	var rt = RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		assert.Same(t, xReq, req)
		return xResp, xErr
	})

	resp, err := rt.RoundTrip(xReq)
	assert.Same(t, xResp, resp)
	assert.Same(t, xErr, err)
}

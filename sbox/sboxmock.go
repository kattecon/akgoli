package sbox

import (
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

type sboxSvcMockImpl struct{}

// NewSBoxSvcMock returns an SBoxSvc for tests that uses URL-query-escaped
// JSON instead of real encryption. The encoded output is deterministic and
// human-readable, making test assertions easier.
//
// Differences from the production SBoxSvc:
//   - Output is deterministic. The same input always produces the same string.
//   - Any mock instance can decode any other mock instance's output.
//   - Authentication and wrong-key failures cannot be exercised.
//   - A URL-unescape error in Decode is silently dropped. The malformed
//     input reaches json.Unmarshal, which returns its own error.
func NewSBoxSvcMock() SBoxSvc {
	return sboxSvcMockImpl{}
}

func (sb sboxSvcMockImpl) Encode(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "could not serialize value")
	}

	return url.QueryEscape(string(data)), nil
}

func (sb sboxSvcMockImpl) Decode(encoded string, value any) error {
	// Unescape error is intentionally dropped. If the input is malformed,
	// json.Unmarshal below reports its own error against the raw string.
	r, _ := url.QueryUnescape(encoded)
	return json.Unmarshal([]byte(r), value)
}

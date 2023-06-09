package sbox

import (
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

type sboxSvcMockImpl struct{}

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
	r, _ := url.QueryUnescape(encoded)
	return json.Unmarshal([]byte(r), value)
}

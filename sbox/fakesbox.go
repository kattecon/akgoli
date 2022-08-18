package sbox

import (
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
)

type FakeSBoxSvcImpl struct{}

func NewFakeSBoxSvc() FakeSBoxSvcImpl {
	return FakeSBoxSvcImpl{}
}

func (sb FakeSBoxSvcImpl) Encode(value any) (string, error) {
	// Serialize value
	data, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "could not serialize value")
	}

	return url.QueryEscape(string(data)), nil
}

func (sb FakeSBoxSvcImpl) Decode(encoded string, value any) error {
	r, _ := url.QueryUnescape(encoded)
	return json.Unmarshal([]byte(r), value)
}

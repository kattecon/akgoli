package sbox

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/akshaal/akgoli/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"
)

// Encrypts and authenticates small messages. Wrapper around "secretbox".
// Unlike secretbox, sbox takes care of nonce and keys itself.

const (
	secretKeySize = 32
	nonceSize     = 24

	ErrFailedToDecrypt  = utils.ConstError("failed to decrypt")
	ErrWrongEncodedSize = utils.ConstError("strange encoded size")
)

type SBoxSvc interface {
	Encode(value any) (string, error)
	Decode(encoded string, value any) error
}

type sboxSvcImpl struct {
	secretKey [secretKeySize]byte
}

func NewSBoxSvc() SBoxSvc {
	var sb sboxSvcImpl

	nr, err := rand.Read(sb.secretKey[:])
	if nr != secretKeySize || err != nil {
		panic(err)
	}

	return &sb
}

func (sb *sboxSvcImpl) Encode(value any) (string, error) {
	// Serialize value
	data, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "could not serialize value")
	}

	// Generate nonce
	var nonce [nonceSize]byte
	nr, err := rand.Read(nonce[:])
	if nr != nonceSize || err != nil {
		panic(err)
	}

	// Encrypt data. First argument to the seal says that nonce must be prepended in the encrypted data
	encrypted := secretbox.Seal(nonce[:], data, &nonce, &sb.secretKey)

	// Encode
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(encrypted)))
	base64.URLEncoding.Encode(encoded, encrypted)

	return string(encoded), nil
}

func (sb *sboxSvcImpl) Decode(encoded string, value any) error {
	// NOTE: No errors are wrapped to avoid extra allocations

	if len(encoded) <= nonceSize {
		return ErrWrongEncodedSize
	}

	encodedBytes := []byte(encoded)
	decodeBytes := make([]byte, base64.URLEncoding.DecodedLen(len(encodedBytes)))
	decodedSize, err := base64.URLEncoding.Decode(decodeBytes, encodedBytes)
	if err != nil {
		return err
	}
	decoded := decodeBytes[:decodedSize]

	var decryptNonce [nonceSize]byte
	copy(decryptNonce[:], decoded[:nonceSize])
	decrypted, valid := secretbox.Open(
		nil,
		decoded[nonceSize:],
		&decryptNonce,
		&sb.secretKey,
	)
	if !valid {
		return ErrFailedToDecrypt
	}

	return json.Unmarshal(decrypted, value)
}

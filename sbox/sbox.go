// Package sbox provides a high-level secure encryption service for small data payloads.
//
// This package wraps the NaCl secretbox cryptographic library to provide authenticated encryption
// with automatic key and nonce management. It's designed for encrypting small data structures
// that can be JSON-serialized, such as session tokens, API keys, or configuration data.
//
// Key Features:
//   - Authenticated encryption using ChaCha20-Poly1305 (via NaCl secretbox).
//   - Automatic cryptographically secure key generation.
//   - Unique random nonce for each encryption operation.
//   - JSON serialization/deserialization of arbitrary Go values.
//   - URL-safe base64 encoding for easy transport.
//   - Protection against tampering and replay attacks.
//
// Security Properties:
//   - Confidentiality: Data is encrypted and cannot be read without the key. Key is never exposed.
//   - Authenticity: Tampering with encrypted data will be detected during decryption.
//   - Semantic security: Identical plaintexts produce different ciphertexts.
//   - Forward secrecy: Each service instance uses a unique ephemeral key.
//
// Example Usage:
//
//	// Create a new encryption service.
//	svc := sbox.NewSBoxSvc()
//
//	// Encrypt any JSON-serializable data.
//	data := map[string]any{
//		"user_id": 12345,
//		"role":    "admin",
//		"exp":     time.Now().Add(time.Hour).Unix(),
//	}
//
//	encrypted, err := svc.Encode(data)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Decrypt back to original data.
//	var decrypted map[string]any
//	err = svc.Decode(encrypted, &decrypted)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Thread Safety:
// SBoxSvc instances are safe for concurrent use by multiple goroutines after creation.
//
// Limitations:
//   - Designed for small payloads (typically < 1MB due to JSON overhead).
//   - Each service instance has a unique key - data encrypted by one instance
//     cannot be decrypted by another instance.
//   - Keys are ephemeral and not persisted - service restart loses all keys.
package sbox

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/kattecon/akgoli/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	// secretKeySize is the size in bytes of the secret key used for encryption.
	// This matches the key size required by NaCl secretbox (32 bytes for ChaCha20).
	secretKeySize = 32

	// nonceSize is the size in bytes of the nonce used for encryption.
	// This matches the nonce size required by NaCl secretbox (24 bytes).
	nonceSize = 24

	// ErrFailedToDecrypt is returned when decryption fails due to authentication
	// failure, wrong key, or corrupted data.
	ErrFailedToDecrypt = utils.ConstError("failed to decrypt")

	// ErrWrongEncodedSize is returned when the encoded string is too short to
	// contain a valid encrypted message (must be at least nonceSize bytes after base64 decoding).
	ErrWrongEncodedSize = utils.ConstError("strange encoded size")
)

// SBoxSvc provides authenticated encryption services for small data payloads.
//
// Each instance maintains its own ephemeral encryption key and can only decrypt
// data that it previously encrypted. This provides strong isolation between
// different service instances.
//
// The service automatically handles:
//   - Key generation using cryptographically secure randomness.
//   - Nonce generation for each encryption operation.
//   - JSON serialization/deserialization.
//   - Base64 encoding for transport.
//   - Authentication tag verification during decryption.
type SBoxSvc interface {
	// Encode encrypts and encodes the given value into a URL-safe base64 string.
	//
	// The value must be JSON-serializable. The resulting string contains:
	//   - A random 24-byte nonce (prepended to ciphertext).
	//   - The encrypted JSON representation of the value.
	//   - An authentication tag to prevent tampering.
	//
	// Each call to Encode with the same value produces a different result due to
	// the random nonce, providing semantic security.
	//
	// Returns an error if JSON serialization fails or if the cryptographic
	// random number generator fails (which would indicate a serious system problem).
	Encode(value any) (string, error)

	// Decode decrypts and decodes an encrypted string back into the provided value.
	//
	// The encoded string must have been created by the same SBoxSvc instance.
	// The value parameter must be a pointer to the type you want to decode into.
	//
	// Returns:
	//   - ErrWrongEncodedSize if the string is too short to be valid.
	//   - ErrFailedToDecrypt if authentication fails or data is corrupted.
	//   - base64 decoding errors for malformed input.
	//   - JSON unmarshaling errors if the decrypted data doesn't match the target type.
	//
	// Example:
	//   var result MyStruct
	//   err := svc.Decode(encrypted, &result)
	Decode(encoded string, value any) error
}

// sboxSvcImpl implements the SBoxSvc interface using NaCl secretbox for encryption.
// Each instance maintains its own unique secret key that is generated during construction.
type sboxSvcImpl struct {
	// secretKey is the 32-byte key used for encryption/decryption.
	// This key is generated once during construction and never changes.
	secretKey [secretKeySize]byte
}

// NewSBoxSvc creates a new SBoxSvc instance with a randomly generated encryption key.
//
// Each instance has its own unique key, so data encrypted by one instance cannot
// be decrypted by another instance. This provides strong isolation between services.
//
// The function panics if the system's cryptographic random number generator fails,
// which would indicate a serious system security problem that should not be ignored.
//
// Returns a new SBoxSvc ready for immediate use.
func NewSBoxSvc() SBoxSvc {
	var sb sboxSvcImpl

	nr, err := rand.Read(sb.secretKey[:])
	if nr != secretKeySize || err != nil {
		panic(err)
	}

	return &sb
}

// Encode encrypts the given value and returns it as a URL-safe base64 encoded string.
//
// The same value encrypted multiple times will produce different results due to
// the random nonce, providing semantic security against pattern analysis.
//
// Security guarantees:
//   - Confidentiality: The original data cannot be recovered without the key.
//   - Authenticity: Any tampering with the result will be detected during decryption.
//   - Freshness: Each encryption uses a unique nonce.
//
// Parameters:
//   - value: Any Go value that can be JSON-marshaled.
//
// Returns:
//   - A URL-safe base64 encoded string containing the encrypted data.
//   - An error if JSON marshaling fails.
//
// Panics if the cryptographic random number generator fails, indicating a
// serious system security issue.
func (sb *sboxSvcImpl) Encode(value any) (string, error) {
	// Serialize value to JSON for encryption.
	data, err := json.Marshal(value)
	if err != nil {
		return "", errors.Wrap(err, "could not serialize value")
	}

	// Generate a cryptographically secure random nonce.
	// Each encryption operation must use a unique nonce for security.
	var nonce [nonceSize]byte
	nr, err := rand.Read(nonce[:])
	if nr != nonceSize || err != nil {
		panic(err)
	}

	// Encrypt the data using NaCl secretbox (ChaCha20-Poly1305).
	// The nonce is prepended to the encrypted data for later extraction during decryption.
	encrypted := secretbox.Seal(nonce[:], data, &nonce, &sb.secretKey)

	// Encode the encrypted data as URL-safe base64 for transport.
	// This format is safe for use in URLs, headers, and JSON without escaping.
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(encrypted)))
	base64.URLEncoding.Encode(encoded, encrypted)

	return string(encoded), nil
}

// Decode decrypts a base64 encoded string and unmarshals it into the provided value.
//
// Security verification:
//   - Authentication tag is verified to detect tampering.
//   - Only data encrypted by this specific service instance can be decrypted.
//   - Malformed or corrupted data will be rejected.
//
// Parameters:
//   - encoded: A base64 string previously returned by Encode.
//   - value: A pointer to the variable where the decrypted data should be stored.
//
// Returns:
//   - ErrWrongEncodedSize if the string is too short to contain valid encrypted data.
//   - ErrFailedToDecrypt if authentication fails (wrong key, corruption, or tampering).
//   - Base64 decoding errors for malformed input.
//   - JSON unmarshaling errors if the decrypted data doesn't match the target type.
//
// Example:
//
//	var myData MyStruct
//	err := svc.Decode(encryptedString, &myData)
func (sb *sboxSvcImpl) Decode(encoded string, value any) error {
	// Validate minimum length: must contain at least a nonce.
	// Note: Errors are not wrapped here to avoid extra allocations in the hot path.
	if len(encoded) <= nonceSize {
		return ErrWrongEncodedSize
	}

	// Decode from URL-safe base64.
	encodedBytes := []byte(encoded)
	decodeBytes := make([]byte, base64.URLEncoding.DecodedLen(len(encodedBytes)))
	decodedSize, err := base64.URLEncoding.Decode(decodeBytes, encodedBytes)
	if err != nil {
		return err
	}
	decoded := decodeBytes[:decodedSize]

	// Extract the nonce from the first 24 bytes.
	var decryptNonce [nonceSize]byte
	copy(decryptNonce[:], decoded[:nonceSize])

	// Decrypt and authenticate the data using NaCl secretbox.
	// This will fail if the data has been tampered with or was encrypted with a different key.
	decrypted, valid := secretbox.Open(
		nil,
		decoded[nonceSize:],
		&decryptNonce,
		&sb.secretKey,
	)
	if !valid {
		return ErrFailedToDecrypt
	}

	// Unmarshal the JSON back into the provided value.
	return json.Unmarshal(decrypted, value)
}

package sbox

import (
	"math"
	"strings"
	"testing"

	"github.com/agnivade/levenshtein"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	A string
	B float64
}

type complexStruct struct {
	Name        string
	Age         int
	Scores      []float64
	Metadata    map[string]any
	IsActive    bool
	Coordinates [2]float64
}

func TestIt(t *testing.T) {
	sb1 := NewSBoxSvc()
	sb2 := NewSBoxSvc()
	assert.NotEqual(t, sb1, sb2)

	_, err := sb1.Encode(testStruct{B: math.NaN()})
	assert.NotNil(t, err)

	msg1 := testStruct{A: "hello-world", B: 32}

	enc1a, err := sb1.Encode(msg1)
	assert.Nil(t, err)

	enc1b, err := sb1.Encode(msg1)
	assert.Nil(t, err)

	assert.NotEqual(t, enc1a, enc1b)

	enc2a, err := sb2.Encode(msg1)
	assert.Nil(t, err)

	enc2b, err := sb2.Encode(msg1)
	assert.Nil(t, err)

	assert.NotEqual(t, enc2a, enc2b)

	assert.NotEqual(t, enc1a, enc2a)
	assert.NotEqual(t, enc1a, enc2b)
	assert.NotEqual(t, enc1b, enc2a)
	assert.NotEqual(t, enc1b, enc2b)

	assert.NotContains(t, enc1a, msg1.A)
	assert.NotContains(t, enc1b, msg1.A)
	assert.NotContains(t, enc2a, msg1.A)
	assert.NotContains(t, enc2b, msg1.A)

	// Decode OK - - -- -

	var r1a testStruct
	assert.Nil(t, sb1.Decode(enc1a, &r1a))
	assert.Equal(t, msg1, r1a)

	var r1b testStruct
	assert.Nil(t, sb1.Decode(enc1b, &r1b))
	assert.Equal(t, msg1, r1b)

	var r2a testStruct
	assert.Nil(t, sb2.Decode(enc2a, &r2a))
	assert.Equal(t, msg1, r2a)

	var r2b testStruct
	assert.Nil(t, sb2.Decode(enc2b, &r2b))
	assert.Equal(t, msg1, r2b)

	// Decode not OK - - -- -

	var r testStruct
	assert.Equal(t, ErrFailedToDecrypt, sb2.Decode(enc1a, &r))
	assert.Equal(t, ErrFailedToDecrypt, sb2.Decode(enc1b, &r))
	assert.Equal(t, ErrFailedToDecrypt, sb1.Decode(enc2a, &r))
	assert.Equal(t, ErrFailedToDecrypt, sb1.Decode(enc2b, &r))
	assert.Equal(t, ErrWrongEncodedSize, sb2.Decode("", &r))
	assert.NotNil(t, sb2.Decode(strings.Repeat(" ", 200), &r))
}

// TestEmptyAndNilValues tests edge cases with empty/nil data.
func TestEmptyAndNilValues(t *testing.T) {
	sb := NewSBoxSvc()

	// Test empty string.
	t.Run("empty_string", func(t *testing.T) {
		original := ""
		encoded, err := sb.Encode(original)
		assert.NoError(t, err)

		var decoded string
		err = sb.Decode(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	// Test empty struct.
	t.Run("empty_struct", func(t *testing.T) {
		original := testStruct{}
		encoded, err := sb.Encode(original)
		assert.NoError(t, err)

		var decoded testStruct
		err = sb.Decode(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	// Test nil pointer.
	t.Run("nil_pointer", func(t *testing.T) {
		var original *testStruct
		encoded, err := sb.Encode(original)
		assert.NoError(t, err)

		var decoded *testStruct
		err = sb.Decode(encoded, &decoded)
		assert.NoError(t, err)
		assert.Nil(t, decoded)
	})
}

// TestLargeData tests behavior with large data structures.
func TestLargeData(t *testing.T) {
	sb := NewSBoxSvc()

	t.Run("large_string", func(t *testing.T) {
		original := strings.Repeat("abcdefghij", 1000)
		encoded, err := sb.Encode(original)
		assert.NoError(t, err)

		var decoded string
		err = sb.Decode(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	t.Run("large_slice", func(t *testing.T) {
		// Test with a large slice.
		original := make([]int, 1000)
		for i := range original {
			original[i] = i * i
		}
		encoded, err := sb.Encode(original)
		assert.NoError(t, err)

		var decoded []int
		err = sb.Decode(encoded, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})
}

// TestSpecialCharacters tests Unicode and special characters.
func TestSpecialCharacters(t *testing.T) {
	sb := NewSBoxSvc()

	testCases := []string{
		"Hello 世界",
		"🚀🌟✨",
		"line1\nline2\ttab",
		"\"quotes\" and 'apostrophes'",
	}

	for _, original := range testCases {
		t.Run("special_chars", func(t *testing.T) {
			encoded, err := sb.Encode(original)
			assert.NoError(t, err)

			var decoded string
			err = sb.Decode(encoded, &decoded)
			assert.NoError(t, err)
			assert.Equal(t, original, decoded)
		})
	}
}

// TestMalformedInput tests various malformed input scenarios.
func TestMalformedInput(t *testing.T) {
	sb := NewSBoxSvc()

	t.Run("various_invalid_lengths", func(t *testing.T) {
		var r testStruct

		for i := 0; i < 24; i++ {
			input := strings.Repeat("a", i)
			err := sb.Decode(input, &r)
			assert.Equal(t, ErrWrongEncodedSize, err, "Length %d should return ErrWrongEncodedSize", i)
		}
	})

	t.Run("invalid_base64_characters", func(t *testing.T) {
		var r testStruct

		// Characters that are not valid base64.
		invalidInputs := []string{
			strings.Repeat("!", 50),  // Invalid characters.
			strings.Repeat("@", 50),  // Invalid characters.
			strings.Repeat("\\", 50), // Invalid characters.
			strings.Repeat(" ", 50),  // Spaces (invalid for URLEncoding).
		}

		for _, input := range invalidInputs {
			err := sb.Decode(input, &r)
			assert.Error(t, err)
			assert.NotEqual(t, ErrWrongEncodedSize, err)
			assert.NotEqual(t, ErrFailedToDecrypt, err)
		}
	})

	t.Run("corrupted_valid_encoding", func(t *testing.T) {
		// Create a valid encoding first.
		original := testStruct{A: "test", B: 123.45}
		encoded, err := sb.Encode(original)
		assert.NoError(t, err)

		// Corrupt it in various ways that should definitely fail.
		corruptedInputs := []string{
			encoded[:len(encoded)-5] + "XXXXX",                                // Change last 5 chars.
			"AAAAA" + encoded[5:],                                             // Change first 5 chars.
			encoded[:len(encoded)/2] + "BADDATA" + encoded[len(encoded)/2+7:], // Corrupt middle.
		}

		for i, corrupted := range corruptedInputs {
			var r testStruct
			err := sb.Decode(corrupted, &r)

			// Should be either ErrFailedToDecrypt or a base64 error, but not success.
			assert.Error(t, err, "Corrupted input %d should fail to decode", i)
		}
	})
}

// TestTypeConversions tests decoding to wrong types.
func TestTypeConversions(t *testing.T) {
	sb := NewSBoxSvc()

	t.Run("string_to_int", func(t *testing.T) {
		// Encode a string, try to decode as int.
		encoded, err := sb.Encode("not a number")
		assert.NoError(t, err)

		var decoded int
		err = sb.Decode(encoded, &decoded)
		assert.Error(t, err) // Should fail.
	})

	t.Run("struct_to_string", func(t *testing.T) {
		// Encode a struct, try to decode as string.
		original := testStruct{A: "test", B: 1.23}
		encoded, err := sb.Encode(original)
		assert.NoError(t, err)

		var decoded string
		err = sb.Decode(encoded, &decoded)
		assert.Error(t, err) // Should fail.
	})
}

// TestMultipleServices tests that different services can't decrypt each other's data.
func TestMultipleServices(t *testing.T) {
	services := []SBoxSvc{
		NewSBoxSvc(),
		NewSBoxSvc(),
		NewSBoxSvc(),
	}

	original := testStruct{A: "secret", B: 42.0}

	// Each service encodes the same data.
	encodings := make([]string, len(services))
	for i, svc := range services {
		var err error
		encodings[i], err = svc.Encode(original)
		assert.NoError(t, err)
	}

	// Verify all encodings are unique.
	for i := 0; i < len(encodings); i++ {
		for j := i + 1; j < len(encodings); j++ {
			assert.NotEqual(t, encodings[i], encodings[j])
		}
	}

	// Each service can only decode its own encoding.
	for i, svc := range services {
		for j, encoding := range encodings {
			var decoded testStruct
			err := svc.Decode(encoding, &decoded)

			if i == j {
				assert.NoError(t, err)
				assert.Equal(t, original, decoded)
			} else {
				assert.Equal(t, ErrFailedToDecrypt, err)
			}
		}
	}
}

// TestConsistentErrors tests that error conditions are consistent.
func TestConsistentErrors(t *testing.T) {
	sb := NewSBoxSvc()

	t.Run("empty_string_always_wrong_size", func(t *testing.T) {
		var r testStruct

		// Test empty string multiple times.
		for i := 0; i < 5; i++ {
			err := sb.Decode("", &r)
			assert.Equal(t, ErrWrongEncodedSize, err)
		}
	})

	t.Run("cross_service_always_decrypt_fail", func(t *testing.T) {
		sb1 := NewSBoxSvc()
		sb2 := NewSBoxSvc()

		encoded, err := sb1.Encode(testStruct{A: "test", B: 1.0})
		assert.NoError(t, err)

		// Try decoding with different service multiple times.
		for i := 0; i < 5; i++ {
			var r testStruct
			err := sb2.Decode(encoded, &r)
			assert.Equal(t, ErrFailedToDecrypt, err)
		}
	})
}

// TestEncodingLevenshteinDistance tests that same data encrypted twice produces sufficiently different results.
func TestEncodingLevenshteinDistance(t *testing.T) {
	sb := NewSBoxSvc()

	// Test with different data types and sizes.
	testCases := []struct {
		name string
		data any
	}{
		{"string", "test-data-for-encryption"},
		{"struct", testStruct{A: "hello-world", B: 123.456}},
		{"large_string", strings.Repeat("data", 50)}, // 200 chars
		{"complex_struct", complexStruct{
			Name:        "Test User",
			Age:         30,
			Scores:      []float64{85.5, 92.0, 78.5},
			Metadata:    map[string]any{"level": "expert", "active": true},
			IsActive:    true,
			Coordinates: [2]float64{40.7128, -74.0060},
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test 1000 times.
			for i := 0; i < 1000; i++ {
				// Encrypt the same data twice.
				enc1, err := sb.Encode(tc.data)
				assert.NoError(t, err)

				enc2, err := sb.Encode(tc.data)
				assert.NoError(t, err)

				// Ensure they're different.
				assert.NotEqual(t, enc1, enc2)

				// Calculate Levenshtein distance as percentage.
				maxLen := len(enc1)
				if len(enc2) > maxLen {
					maxLen = len(enc2)
				}

				dist := float64(levenshtein.ComputeDistance(enc1, enc2)) / float64(maxLen)

				assert.Greater(t, dist, 0.8, "Levenshtein distance %.3f should be > 0.8 for encryptions %d: %s vs %s", dist, i, enc1, enc2)

				// Verify both can be decrypted correctly.
				var decoded1, decoded2 any
				switch tc.data.(type) {
				case string:
					var d1, d2 string
					assert.NoError(t, sb.Decode(enc1, &d1))
					assert.NoError(t, sb.Decode(enc2, &d2))
					decoded1, decoded2 = d1, d2
				case testStruct:
					var d1, d2 testStruct
					assert.NoError(t, sb.Decode(enc1, &d1))
					assert.NoError(t, sb.Decode(enc2, &d2))
					decoded1, decoded2 = d1, d2
				case complexStruct:
					var d1, d2 complexStruct
					assert.NoError(t, sb.Decode(enc1, &d1))
					assert.NoError(t, sb.Decode(enc2, &d2))
					decoded1, decoded2 = d1, d2
				}

				assert.Equal(t, tc.data, decoded1)
				assert.Equal(t, tc.data, decoded2)
			}
		})
	}
}

// TestDataIntegrity verifies that encrypted data doesn't leak information.
func TestDataIntegrity(t *testing.T) {
	sb := NewSBoxSvc()

	t.Run("no_plaintext_leakage", func(t *testing.T) {
		secrets := []string{
			"password123",
			"secret-api-key",
			"confidential-data",
			"social-security-number",
		}

		for _, secret := range secrets {
			encoded, err := sb.Encode(secret)
			assert.NoError(t, err)

			// The secret should not appear in the encoded form.
			assert.NotContains(t, encoded, secret)
			assert.NotContains(t, strings.ToLower(encoded), strings.ToLower(secret))
		}
	})
}

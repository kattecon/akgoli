package sbox

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	A string
	B float64
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

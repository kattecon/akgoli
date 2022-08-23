package typesafe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncMap(t *testing.T) {
	var m SyncMap[string, int]

	_, loaded := m.Load("x")
	assert.False(t, loaded)

	m.Store("x", 123)

	m.Store("y", 999)
	m.Delete("y")

	v, loaded := m.Load("x")
	assert.True(t, loaded)
	assert.Equal(t, 123, v)

	v, loaded = m.LoadOrStore("x", 321)
	assert.True(t, loaded)
	assert.Equal(t, 123, v)

	_, loaded = m.LoadOrStore("y", 444)
	assert.False(t, loaded)

	v, loaded = m.LoadOrStore("y", 555)
	assert.True(t, loaded)
	assert.Equal(t, 444, v)

	v, loaded = m.LoadOrStore("y", 333)
	assert.True(t, loaded)
	assert.Equal(t, 444, v)

	_, loaded = m.LoadAndDelete("z")
	assert.False(t, loaded)

	m.Store("z", 5311)

	v, loaded = m.LoadAndDelete("z")
	assert.True(t, loaded)
	assert.Equal(t, 5311, v)

	_, loaded = m.Load("z")
	assert.False(t, loaded)

	var xFound, yFound, otherFound bool
	m.Range(func(key string, value int) bool {
		if key == "x" && value == 123 {
			xFound = true
		} else if key == "y" && value == 444 {
			yFound = true
		} else {
			otherFound = true
		}
		return true
	})
	assert.True(t, xFound)
	assert.True(t, yFound)
	assert.False(t, otherFound)
}

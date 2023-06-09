package utils

import (
	"testing"

	"github.com/kattecon/akgoli/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestMasked(t *testing.T) {
	buflog := testutils.NewBufferingLogger(zap.InfoLevel)
	buflog.Logger.Info("test", Masked("aaa", "abcd"))
	assert.Equal(t, "{'level':'info','msg':'test','aaa':'****'}\n", buflog.JsonNoDoubleQuotes())
}

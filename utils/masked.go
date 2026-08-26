package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type masked struct {
	k string
	v string
}

// Masked returns a zap.Field that logs the value v under key k with every
// rune replaced by an asterisk. The rune count of the original value is
// preserved and visible in the log output. Use this for secrets that must
// appear in structured logs without exposing their content.
// See also MaskAll for raw string masking without zap.
func Masked(k string, v string) zap.Field {
	return zap.Inline(&masked{k, v})
}

func (m *masked) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString(m.k, MaskAll(m.v))
	return nil
}

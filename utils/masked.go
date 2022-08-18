package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type masked struct {
	k string
	v string
}

func Masked(k string, v string) zap.Field {
	return zap.Inline(&masked{k, v})
}

func (m *masked) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString(m.k, MaskAll(m.v))
	return nil
}

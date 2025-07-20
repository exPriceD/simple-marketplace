package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/exPriceD/simple-marketplace/internal/platform/log"
)

type ZapLogger struct {
	l *zap.Logger
}

func New(level string) (*ZapLogger, error) {
	var lvl zapcore.Level
	switch level {
	case "debug":
		lvl = zap.DebugLevel
	case "info", "":
		lvl = zap.InfoLevel
	default:
		lvl = zap.InfoLevel
	}
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(lvl),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:       "ts",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stack",
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime:    zapcore.RFC3339TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	z, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{l: z}, nil
}

func (z *ZapLogger) Info(_ context.Context, msg string, kv ...any) {
	z.l.Info(msg, fields(kv...)...)
}

func (z *ZapLogger) Error(_ context.Context, msg string, kv ...any) {
	z.l.Error(msg, fields(kv...)...)
}

func (z *ZapLogger) Debug(_ context.Context, msg string, kv ...any) {
	z.l.Debug(msg, fields(kv...)...)
}

func (z *ZapLogger) With(kv ...any) log.Logger {
	return &ZapLogger{l: z.l.With(fields(kv...)...)}
}

func fields(kv ...any) []zap.Field {
	fs := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		key, ok := kv[i].(string)
		if !ok {
			continue
		}
		fs = append(fs, zap.Any(key, kv[i+1]))
	}
	return fs
}

var _ log.Logger = (*ZapLogger)(nil)

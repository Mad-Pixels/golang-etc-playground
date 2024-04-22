package log

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	*zap.Logger
}

func New(level []byte) (*Log, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.DisableCaller = true

	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText(level); err != nil {
		return nil, fmt.Errorf("failed to initialize zap log level from '%s': %s", level, err)
	}
	cfg.Level = zap.NewAtomicLevelAt(logLevel)
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := cfg.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return nil, err
	}
	return &Log{
		Logger: logger.WithOptions(zap.AddStacktrace(zap.DPanicLevel)),
	}, nil
}

// Printf implementation.
func (l *Log) Printf(format string, args ...interface{}) {
	var fields []zap.Field
	var etc []interface{}

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]interface{}:
			for k, val := range v {
				fields = append(fields, zap.Any(k, val))
			}
		default:
			etc = append(etc, v)
		}
	}
	msg := fmt.Sprintf(format, etc...)

	if strings.Contains(strings.ToLower(msg), "error") {
		l.Logger.Error(msg, fields...)
	} else {
		l.Logger.Info(msg, fields...)
	}
}

package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env, level string) (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	if env == "prod" {
		cfg = zap.NewProductionConfig()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)

	return cfg.Build()
}

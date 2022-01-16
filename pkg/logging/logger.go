package logging

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Dev() *zap.Logger {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapConfig.EncoderConfig.EncodeTime = nil
	zapConfig.DisableStacktrace = true
	lg, err := zapConfig.Build()
	if err != nil {
		log.Fatal("failed to initialize logging")
	}

	return lg
}

// TODO production logger
func Prod() *zap.Logger {
	return nil
}

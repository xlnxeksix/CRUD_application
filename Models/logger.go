package models

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	// Change the log output to a file called "project_log.txt"
	config.OutputPaths = []string{"/Logs/project_log.txt"}
	config.ErrorOutputPaths = []string{"/Logs/project_log.txt"}
	// Customize other logger options as needed (e.g., log level, encoding, etc.)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// Create the logger
	logger, err := config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	Logger = logger
}

func CloseLogger() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

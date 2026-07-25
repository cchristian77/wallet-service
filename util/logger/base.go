package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/cchristian77/wallet-service/util/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Initialize() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.FunctionKey = "func"

	// log to console
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	cores := []zapcore.Core{consoleCore}

	// log to file (if found)
	_ = os.MkdirAll("logs", 0755)
	logFile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(logFile),
			zapcore.InfoLevel,
		)
		cores = append(cores, fileCore)
	}

	logger = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger
}

// InitializeNop sets a no-op logger for unit tests.
func InitializeNop() {
	logger = zap.NewNop()
}

// Get returns the global logger instance.
func Get() *zap.Logger {
	if logger == nil {
		return Initialize()
	}
	return logger
}

// L call logger instance without the need for context
func L() *zap.Logger {
	return Get().WithOptions(zap.AddCallerSkip(1))
}

func Fatal(ctx context.Context, message string, args ...any) {
	Get().Fatal(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func Error(ctx context.Context, message string, args ...any) {
	Get().Error(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func Warn(ctx context.Context, message string, args ...any) {
	Get().Warn(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func Info(ctx context.Context, message string, args ...any) {
	Get().Info(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func Debug(ctx context.Context, message string, args ...any) {
	Get().Debug(fmt.Sprintf(message, args...), getZapFieldsFromCtx(ctx)...)
}

func getZapFieldsFromCtx(ctx context.Context) []zapcore.Field {
	if ctx == nil {
		return nil
	}

	var fields []zapcore.Field

	correlationID := constant.CorrelationIDFromCtx(ctx)
	if correlationID != "" {
		fields = append(fields, zap.String(constant.XCorrelationIDKey, correlationID))
	}

	return fields
}

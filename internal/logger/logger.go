package logger

import (
	"os"
	"time"

	"github.com/mathif92/auth-service/internal/logger/encoders"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger(level zapcore.Level) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     rfc3399NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	encoder := encoders.NewKeyValueEncoder(encoderConfig)

	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(os.Stderr)), level))
}

// rfc3399NanoTimeEncoder serializes a time.Time to an RFC3399-formatted string
// with microsecond precision padded with zeroes to make it fixed width.
func rfc3399NanoTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"

	enc.AppendString(t.UTC().Format(RFC3339Micro))
}

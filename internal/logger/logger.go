package logger

import (
    "os"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(level string) error {
    logLevel := zapcore.InfoLevel
    if level == "debug" {
        logLevel = zapcore.DebugLevel
    }

    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.CapitalColorLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.StringDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    core := zapcore.NewCore(
        zapcore.NewConsoleEncoder(encoderConfig),
        zapcore.AddSync(os.Stdout),
        logLevel,
    )

    Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
    return nil
}

func Sync() {
    if Log != nil {
        _ = Log.Sync()
    }
}
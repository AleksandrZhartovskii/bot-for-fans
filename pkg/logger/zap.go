package logger

import (
    "fmt"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "os"
    "strings"
)

const (
    LogLevelDebug   = "DEBUG"
    LogLevelInfo    = "INFO"
    LogLevelWarning = "WARN"
    LogLevelError   = "ERROR"

	currentLogsPath = "_logs"
)

func NewLogger(serviceName, logLevel string) (*zap.Logger, error) {
    // Create directory _logs if not exist
    if _, err := os.Stat(currentLogsPath); os.IsNotExist(err) {
        err = os.Mkdir(currentLogsPath, 0777)
        if err != nil {
            return nil, err
        }
    }

    // Open log file by path _logs/{serviceName}.log
    logFilePath := fmt.Sprintf("%s/%s.log", currentLogsPath, serviceName)
    _, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
    if err != nil {
        return nil, err
    }

    // Determining the logging level
    var level zapcore.Level
    switch strings.ToUpper(logLevel) {
    case LogLevelDebug:
        level = zap.DebugLevel
    case LogLevelInfo:
        level = zap.InfoLevel
    case LogLevelWarning:
        level = zap.WarnLevel
    case LogLevelError:
        level = zap.ErrorLevel
    default:
        level = zap.DebugLevel
    }

    // Configuring logging parameters
    cfg := zap.Config{
        Level:    zap.NewAtomicLevelAt(level),
        Sampling: nil,
        Encoding: "json",
        EncoderConfig: zapcore.EncoderConfig{
            TimeKey:        "timestamp",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "caller",
            FunctionKey:    zapcore.OmitKey,
            MessageKey:     "message",
            StacktraceKey:  "stacktrace",
            LineEnding:     zapcore.DefaultLineEnding,
            EncodeLevel:    zapcore.LowercaseLevelEncoder,
            EncodeTime:     zapcore.ISO8601TimeEncoder,
            EncodeDuration: zapcore.SecondsDurationEncoder,
            EncodeCaller:   zapcore.ShortCallerEncoder,
        },
        OutputPaths:      []string{"stdout", logFilePath},
        ErrorOutputPaths: []string{"stdout"},
        InitialFields:    nil,
    }

    // Build logger
    logger, err := cfg.Build()
    if err != nil {
        return nil, err
    }
    return logger, nil
}

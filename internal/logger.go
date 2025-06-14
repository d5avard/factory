package internal

import (
	"log"
	"net/http"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()
}

func LogRequest(r *http.Request, message string, code int) {
	log.Printf("IP: %s, Method: %s, Path: %s, Status: %d, Size: %d, Message: %s",
		hostname, r.Method, r.URL.Path, code, r.ContentLength, message)
}

type Logger struct {
	zap.Logger
}

func (l *Logger) Close() {
	if err := l.Sync(); err != nil {
		log.Printf("failed to sync logger: %v", err)
	}
}

func (l *Logger) SetLevel(level zapcore.Level) {
	// Set the logging level dynamically
	l.Logger.Core().Enabled(level)
	atomicLevel := zap.NewAtomicLevelAt(level)
	l.Logger = *l.Logger.WithOptions(zap.IncreaseLevel(atomicLevel))
}

func NewLogger(dir, file string) *Logger {
	// Dynamic level that can be changed at runtime
	atomicLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// Log rotation setup using lumberjack
	lumberjackLogger := &lumberjack.Logger{
		Filename:   path.Join(dir, file),
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     28,   // days
		Compress:   true, // compress old files
	}

	// Encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Outputs: file (JSON) and console (human-readable)
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(lumberjackLogger),
		atomicLevel,
	)

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		atomicLevel,
	)

	core := zapcore.NewTee(fileCore, consoleCore)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{*zapLogger}
}

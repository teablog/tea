package logger

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Level zapcore.Level

var (
	Green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	White   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	Yellow  = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	Red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	Blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	Magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	Cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	Reset   = string([]byte{27, 91, 48, 109})
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
	_minLevel = DebugLevel
	_maxLevel = FatalLevel
)

var (
	log   *zap.SugaredLogger
	level = zapcore.InfoLevel
)

func NewDefaultLogger(out io.Writer) {
	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})
	w := zapcore.AddSync(out)
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(config), w, priority)
	logger := zap.New(core)
	log = logger.Sugar()
}

func NewLogger(out io.Writer, l string) {
	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level
	})
	w := zapcore.AddSync(out)
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(config), w, priority)
	logger := zap.New(core)
	log = logger.Sugar()
	SetLevel(l)
}

func SetLevel(l string) {
	switch l {
	case "debug":
		level = zapcore.Level(DebugLevel)
	case "info":
		level = zapcore.Level(InfoLevel)
	case "error":
		level = zapcore.Level(ErrorLevel)
	case "Fatal":
		level = zapcore.Level(FatalLevel)
	}
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	log.Debugf(template, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

func Wrapf(err error, template string, args ...interface{}) {
	buf := new(bytes.Buffer) // the returned data
	_, _ = fmt.Fprintf(buf, template, args...)
	e := errors.WithStack(err.(error))
	_, _ = fmt.Fprintf(buf, "%+v", e)
	Errorf(buf.String())
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
	os.Exit(1)
}

func Fatalf(template string, args ...interface{}) {
	log.Fatalf(template, args...)
	os.Exit(1)
}

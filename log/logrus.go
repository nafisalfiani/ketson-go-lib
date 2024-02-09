package log

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

var logrusOnce = sync.Once{}

type logrusLogger struct {
	log *logrus.Logger
}

func initLogrus(cfg Config) Interface {
	var logger *logrus.Logger
	logrusOnce.Do(func() {
		logger := logrus.New()
		logger.Out = os.Stdout
		logger.ReportCaller = true
		logger.Formatter = &logrus.JSONFormatter{
			PrettyPrint:     true,
			TimestampFormat: time.RFC3339,
			// CallerPrettyfier: caller(),
			// FieldMap: logrus.FieldMap{
			// 	logrus.FieldKeyFile: "caller",
			// },
		}

		level, err := logrus.ParseLevel(cfg.Level)
		if err != nil {
			log.Fatal().Msg(fmt.Sprintf("failed to parse error level with err: %v", err))
		}
		logger.Level = level
	})

	return &logrusLogger{log: logger}
}

func (l *logrusLogger) Trace(ctx context.Context, obj any) {
	l.log.WithFields(getContextFields(ctx)).
		Trace(fmt.Sprint(getCaller(obj)))
}

func (l *logrusLogger) Debug(ctx context.Context, obj any) {
	l.log.WithFields(getContextFields(ctx)).
		Debug(fmt.Sprint(getCaller(obj)))
}

func (l *logrusLogger) Info(ctx context.Context, obj any) {
	l.log.WithFields(getContextFields(ctx)).
		Info(fmt.Sprint(getCaller(obj)))
}

func (l *logrusLogger) Warn(ctx context.Context, obj any) {
	l.log.WithFields(getContextFields(ctx)).
		Warn(fmt.Sprint(getCaller(obj)))
}

func (l *logrusLogger) Error(ctx context.Context, obj any) {
	l.log.WithFields(getContextFields(ctx)).
		Error(fmt.Sprint(getCaller(obj)))
}

func (l *logrusLogger) Fatal(ctx context.Context, obj any) {
	l.log.WithFields(getContextFields(ctx)).
		Fatal(fmt.Sprint(getCaller(obj)))
}

func (l *logrusLogger) Panic(obj any) {
	defer func() { recover() }()
	l.log.WithFields(getPanicStacktrace()).
		Panic(fmt.Sprint(getCaller(obj)))
}

// func caller() func(*runtime.Frame) (function string, file string) {
// 	return func(f *runtime.Frame) (function string, file string) {
// 		p, _ := os.Getwd()

// 		return "", fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
// 	}
// }

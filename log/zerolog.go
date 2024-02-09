package log

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var once = sync.Once{}

type zeroLogger struct {
	log zerolog.Logger
}

func DefaultLogger() Interface {
	return &zeroLogger{
		log: zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3).
			Logger().
			Level(zerolog.DebugLevel),
	}
}

func initZerolog(cfg Config) Interface {
	var zeroLogging zerolog.Logger
	once.Do(func() {
		level, err := zerolog.ParseLevel(cfg.Level)
		if err != nil {
			log.Fatal().Msg(fmt.Sprintf("failed to parse error level with err: %v", err))
		}

		zeroLogging = zerolog.New(os.Stdout).
			With().
			Timestamp().
			CallerWithSkipFrameCount(3).
			Logger().
			Level(level)
	})

	return &zeroLogger{log: zeroLogging}
}

func (l *zeroLogger) Trace(ctx context.Context, obj any) {
	l.log.Trace().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *zeroLogger) Debug(ctx context.Context, obj any) {
	l.log.Debug().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *zeroLogger) Info(ctx context.Context, obj any) {
	l.log.Info().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *zeroLogger) Warn(ctx context.Context, obj any) {
	l.log.Warn().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *zeroLogger) Error(ctx context.Context, obj any) {
	l.log.Error().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *zeroLogger) Fatal(ctx context.Context, obj any) {
	l.log.Fatal().
		Fields(getContextFields(ctx)).
		Msg(fmt.Sprint(getCaller(obj)))
}

func (l *zeroLogger) Panic(obj any) {
	defer func() { recover() }()
	l.log.Panic().
		Fields(getPanicStacktrace()).
		Msg(fmt.Sprint(getCaller(obj)))
}

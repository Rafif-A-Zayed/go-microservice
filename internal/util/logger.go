package util

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"os"
)

func CreateLogger(servicemen string) log.Logger {
	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(logger,
		"svc", servicemen,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
	return logger
}

func Info(logger log.Logger, keysvalue ...interface{}) {

	err := level.Info(logger).Log(keysvalue)
	if err != nil {
		return
	}
}

func Error(logger log.Logger, keysvalue ...interface{}) {

	err := level.Error(logger).Log(keysvalue)
	if err != nil {
		return
	}
}

func Debug(logger log.Logger, keysvalue ...interface{}) {

	err := level.Debug(logger).Log(keysvalue)
	if err != nil {
		return
	}
}

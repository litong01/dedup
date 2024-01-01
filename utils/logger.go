package utils

import (
	"os"

	goslog "golang.org/x/exp/slog"
)

var (
	logger *goslog.Logger
	opts   goslog.HandlerOptions
)

func GetLogger() *goslog.Logger {
	doLog := os.Getenv("DOLOG")
	// TODO getting configuration parameters of the control,
	// then use these parameters to customize the logger.
	if doLog == "" {
		opts.Level = goslog.LevelError
	} else {
		opts.Level = goslog.LevelInfo
	}
	logger = goslog.New(goslog.NewJSONHandler(os.Stdout, &opts))
	goslog.SetDefault(logger)
	return logger
}

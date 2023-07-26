package log

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const loggerKey LoggerKeyType = "loggerKey"

type LoggerKeyType string

func InitLogger(level, component, env string) error {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.LevelFieldName = "log_level" // to compatibility with Graylog

	levelParsed, err := zerolog.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("err parse log level %s: %w", level, err)
	}
	zerolog.SetGlobalLevel(levelParsed)
	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("component", component).
		Logger()

	if env != "" {
		logger = logger.With().Str("env", env).Logger()
	}

	zerolog.DefaultContextLogger = &logger
	return nil
}

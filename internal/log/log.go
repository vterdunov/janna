package log

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	log zerolog.Logger
}

func NewLogger() Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger = logger.With().Str("app", "janna").Logger()
	logger = logger.With().CallerWithSkipFrameCount(3).Logger()
	zerolog.MessageFieldName = "msg"

	return Logger{
		log: logger,
	}
}

func (l *Logger) WithFields(keyvals ...interface{}) Logger {
	fields := make(map[string]interface{})

	for i := 0; i < len(keyvals); i += 2 {
		k := keyvals[i]
		var v interface{} = "(MISSING)"
		if i+1 < len(keyvals) {
			v = keyvals[i+1]
		}

		if keyStr, ok := k.(string); ok {
			if valStr, ok := v.(string); ok {
				fields[keyStr] = valStr
			}
		}
	}

	log := l.log.With().Fields(fields).Logger()
	newLogger := Logger{
		log: log,
	}

	return newLogger
}

func (l *Logger) Info(i interface{}) {
	if str, ok := i.(string); ok {
		l.log.Info().Msg(str)
	}
}

func (l *Logger) Error(err error, i interface{}) {
	var s string
	if str, ok := i.(string); ok {
		s = str
	}

	l.log.Error().Stack().Err(err).Msg(s)
}

package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Config - struct to hold all configurable values
type Config struct {
	AccessToken       string   `env:"ACCESS_TOKEN" envDefault:""`
	LogOutput         string   `env:"LOG_OUTPUT" envDefault:"console"`
	LogLevel          string   `env:"LOG_LEVEL" envDefault:"debug"`
	TimeZone          string   `env:"TIMEZONE" envDefault:"UTC"`
	EscalationPlocies []string `env:"ESCALATION_POLICIES" envDefault:"" envSeparator:","`
}

// Logger - returns a logger configured using the parameters in the struct
func (c Config) Logger() zerolog.Logger {
	zerolog.SetGlobalLevel(strToZerologLevel(c.LogLevel))
	logger := zerolog.New(os.Stderr).With().
		Timestamp().
		Logger()

	if strings.ToLower(c.LogOutput) == "console" {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return logger
}

func strToZerologLevel(level string) zerolog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zerolog.DebugLevel
	case "INFO":
		return zerolog.InfoLevel
	case "WARN":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "FATAL":
		return zerolog.FatalLevel
	case "PANIC":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

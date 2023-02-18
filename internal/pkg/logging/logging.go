package logging

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Initialize takes an io writer to output (defaulting to os.Stderr when nil)
// and a debug option to set the log level to debug.
func Initialize(wr io.Writer, isDebug bool) {
	if wr == nil {
		wr = os.Stderr
	}

	log.SetOutput(wr)
	log.SetFormatter(&logTextFormatter{})

	level := "info"
	if isDebug {
		level = "debug"
	}
	SetLevel(level)
}

// SetLevel parses a string and attempts to set the logging level
// using that string. If the string is not valid the info log level will
// be used as a fall back.
func SetLevel(level string) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		lvl = log.InfoLevel
		log.Warnf("failed to parse log-level '%s', defaulting to 'info'", level)
	}
	log.SetLevel(lvl)
}

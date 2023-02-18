package logging

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type logTextFormatter log.TextFormatter

// Set the logging format to output colored log text
func (f *logTextFormatter) Format(entry *log.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case log.DebugLevel, log.TraceLevel:
		levelColor = 30 // gray
	case log.WarnLevel:
		levelColor = 33 // yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	return []byte(fmt.Sprintf("\x1b[%dm%s\x1b[0m\n", levelColor, entry.Message)), nil
}

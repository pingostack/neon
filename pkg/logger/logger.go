package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var std Logger

type Fields logrus.Fields
type Logger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
	fields logrus.Fields
}

func NewLogger(file string, fields logrus.Fields) *Logger {

	l := &Logger{
		logger: logrus.New(),
	}

	if fields != nil {
		l.SetDefaultFields(fields)
	}

	l.logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	l.logger.Out = os.Stdout
	fd, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		l.entry.Info("Failed to log to file, using default stderr")
	}

	l.logger.SetOutput(io.MultiWriter(fd, os.Stdout))

	return l
}

func (l *Logger) SetDefaultFields(fields logrus.Fields) {
	l.fields = fields

	l.entry = l.logger.WithFields(fields)
}

func (l *Logger) AddDefaultField(key string, value interface{}) {
	if l.fields == nil {
		l.fields = make(map[string]interface{})
	}

	l.fields[key] = value

	l.entry = l.logger.WithFields(l.fields)
}

func (l *Logger) Entry() *logrus.Entry {
	return l.entry
}

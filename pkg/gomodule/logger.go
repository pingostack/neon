package gomodule

import (
	"context"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Fields logrus.Fields

var loggerInstance *loggerModule

type loggerSettings struct {
	Formatter    string `mapstructure:"formatter"`
	Format       string `mapstructure:"format"`
	File         string `mapstructure:"file"`
	Console      bool   `mapstructure:"console"`
	Level        string `mapstructure:"level"`
	ReportCaller bool   `mapstructure:"reportCaller"`
}

type loggerModule struct {
	settings *loggerSettings
}

func init() {
	loggerInstance = &loggerModule{
		settings: &loggerSettings{},
	}
}

func LoggerModule() IModule {
	return loggerInstance
}

func (l *loggerModule) InitModule(ctx context.Context, wg *sync.WaitGroup) (interface{}, error) {

	return l.settings, nil
}

func (l *loggerModule) InitCommand() ([]*cobra.Command, error) {
	// GetRootCmd().PersistentFlags().StringVarP(&l.settings.File, "log", "l", "logs/neon.log", "log file")
	// GetRootCmd().PersistentFlags().StringVarP(&l.settings.Format, "format", "f", "text", "the format of logger")
	// GetRootCmd().PersistentFlags().BoolVar(&l.settings.Console, "console", false, "logger enable console output")

	return nil, nil
}

func (l *loggerModule) ConfigChanged() {

}

func (l *loggerModule) RootCommand(cmd *cobra.Command, args []string) {
	if strings.EqualFold(l.settings.Formatter, "text") {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			ForceColors:     true,
			TimestampFormat: l.settings.Format,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: l.settings.Format,
		})
	}

	logrus.SetReportCaller(l.settings.ReportCaller)
	level, err := logrus.ParseLevel(l.settings.Level)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(level)

	fd, err := createLogFile(l.settings.File)
	if err != nil {
		panic(err)
	}

	var output io.Writer
	if l.settings.Console {
		output = io.MultiWriter(fd, os.Stdout)
	} else {
		output = io.MultiWriter(fd)
	}

	logrus.SetOutput(output)
}

func createLogFile(file string) (*os.File, error) {
	dir := path.Dir(file)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	fd, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return fd, nil
}

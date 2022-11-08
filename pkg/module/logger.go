package module

import (
	"io"
	"os"
	"strings"

	"github.com/pingopenstack/neon/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Fields logrus.Fields

var instance *LoggerModule

type LoggerSettings struct {
	Formatter    string `mapstructure:"formatter"`
	Format       string `mapstructure:"format"`
	File         string `mapstructure:"file"`
	Console      bool   `mapstructure:"console"`
	Level        string `mapstructure:"level"`
	ReportCaller bool   `mapstructure:"reportCaller"`
}

type LoggerModule struct {
	settings *LoggerSettings
}

func init() {
	instance = &LoggerModule{
		settings: &LoggerSettings{},
	}
}

func LoggerModuleInstance() IModule {
	return instance
}

func (l *LoggerModule) OnInitModule() (interface{}, error) {
	return l.settings, nil
}

func (l *LoggerModule) OnInitCommand() ([]*cobra.Command, error) {
	// GetRootCmd().PersistentFlags().StringVarP(&l.settings.File, "log", "l", "logs/neon.log", "log file")
	// GetRootCmd().PersistentFlags().StringVarP(&l.settings.Format, "format", "f", "text", "the format of logger")
	// GetRootCmd().PersistentFlags().BoolVar(&l.settings.Console, "console", false, "logger enable console output")

	return nil, nil
}

func (l *LoggerModule) OnConfigModified() {

}

func (l *LoggerModule) OnPostInitCommand() {
	if strings.ToLower(l.settings.Formatter) == strings.ToLower("text") {
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

	fd, err := utils.CreateDirFile(l.settings.File)
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

func (l *LoggerModule) OnMainRun(cmd *cobra.Command, args []string) {

}

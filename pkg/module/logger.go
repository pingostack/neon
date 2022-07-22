package module

import (
	"io"
	"os"

	"github.com/pingopenstack/neon/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Fields logrus.Fields

var lm *LoggerModule

type LoggerSettings struct {
	Format  string `env:"format"  flag:"format"`
	File    string `env:"file"    flag:"file"`
	Console bool   `env:"console" flag:"console"`
	Level   string `env:"level"   flag:"level"`
}

type LoggerModule struct {
	settings *LoggerSettings
}

func init() {
	lm = &LoggerModule{
		settings: &LoggerSettings{},
	}
}

func LoggerModuleInstance() IModule {
	return lm
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
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	//	logrus.StandardLogger().Out = os.Stdout

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

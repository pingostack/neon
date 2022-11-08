package rtsp

import (
	"github.com/pingopenstack/neon/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RtspFlags struct {
}

type RtspSettings struct {
	RtspServer RtspServerSettings `mapstructure:"server"`
}

type RtspModule struct {
	flags    *RtspFlags
	settings *RtspSettings
	server   *RtspServer
	logger   *logrus.Entry
}

var instance *RtspModule

func init() {
	instance = &RtspModule{
		flags:    &RtspFlags{},
		settings: &RtspSettings{},
	}
	module.Register(instance)
}

func RtspModuleInstance() module.IModule {
	return instance
}

func (w *RtspModule) OnInitModule() (interface{}, error) {
	return w.settings, nil
}

func (w *RtspModule) OnInitCommand() ([]*cobra.Command, error) {
	return nil, nil
}

func (w *RtspModule) OnConfigModified() {

}

func (w *RtspModule) OnPostInitCommand() {
}

func (w *RtspModule) OnMainRun(cmd *cobra.Command, args []string) {
	w.logger = logrus.WithFields(logrus.Fields{
		"Module": "RtspModule",
	})

	var err error
	Logger().Infof("rtsp server listening on %v", w.settings.RtspServer.Addr)

	go func() {
		w.server, err = NewRtspServer(w.settings.RtspServer, w.logger)
		if err != nil {
			Logger().Errorf("New rtsp server failed: %v", err)
			panic(err)
		}
	}()

	Logger().Infof("New rtsp server started")
}

func Logger() *logrus.Entry {
	return instance.logger
}

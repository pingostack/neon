package rtsp

import (
	"github.com/let-light/neon/pkg/module"
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

type RtspContext struct{}

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

func (rm *RtspModule) OnInitModule() (interface{}, error) {
	return rm.settings, nil
}

func (rm *RtspModule) OnInitCommand() ([]*cobra.Command, error) {
	return nil, nil
}

func (rm *RtspModule) OnConfigModified() {

}

func (rm *RtspModule) OnPostInitCommand() {
}

func (rm *RtspModule) OnMainRun(cmd *cobra.Command, args []string) {
	rm.logger = logrus.WithFields(logrus.Fields{
		"Module": "RtspModule",
	})

	var err error
	logger().Infof("rtsp server listening on %v", rm.settings.RtspServer.Addr)

	go func() {
		rm.server, err = NewRtspServer(rm.settings.RtspServer)

		if err != nil {
			logger().Errorf("New rtsp server failed: %v", err)
			panic(err)
		}
	}()

	logger().Infof("New rtsp server started")
}

func rtspContext(s *Session) *RtspContext {
	ctx := s.GetContext(instance)
	if ctx == nil {
		return nil
	}

	return ctx.(*RtspContext)
}

func logger() *logrus.Entry {
	return instance.logger
}

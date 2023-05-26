package sfu

import (
	"context"

	"github.com/let-light/gomodule"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sfuModule *sfu

type ISignalServer interface {
	OnClose(f func())
	Run() error
	Shutdown() error
}

type SfuSettings struct {
	ListenAddr string `json:"listenAddr" mapstructure:"listenAddr"`
	Cert       string `json:"cert" mapstructure:"cert"`
	Key        string `json:"key" mapstructure:"key"`
}

type sfu struct {
	gomodule.DefaultModule
	ctx         context.Context
	preSettings SfuSettings
	settings    *SfuSettings
	ss          *SignalServer
	l           *logrus.Entry
}

func init() {
	sfuModule = &sfu{
		l: logrus.WithField("module", "sfu"),
	}
}

func logger() *logrus.Entry {
	return sfuModule.l
}

func SfuModule() *sfu {
	return sfuModule
}

func (sfu *sfu) InitModule(ctx context.Context, _ *gomodule.Manager) (interface{}, error) {
	sfu.ctx = ctx
	return &sfu.preSettings, nil
}

func (sfu *sfu) InitCommand() ([]*cobra.Command, error) {

	return nil, nil
}

func (sfu *sfu) ConfigChanged() {
	if sfu.settings == nil {
		sfu.settings = &SfuSettings{}
	}

	if *sfu.settings != sfu.preSettings {
		*sfu.settings = sfu.preSettings
	}
}

func (sfu *sfu) ModuleRun() {
	sfu.ss = NewSignalServer(sfu.ctx, SignalServerSettings{
		listenAddr: sfu.settings.ListenAddr,
		cert:       sfu.settings.Cert,
		key:        sfu.settings.Key,
	})

	if e := sfu.ss.Run(); e != nil {
		logger().Errorf("signal server run failed: %v", e)
		panic(e)
	}
}

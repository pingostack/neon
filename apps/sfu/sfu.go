package sfu

import (
	"context"
	"sync"

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
	wg       *sync.WaitGroup
	ctx      context.Context
	settings SfuSettings
	ss       *SignalServer
	l        *logrus.Entry
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

func (sfu *sfu) InitModule(ctx context.Context, wg *sync.WaitGroup) (interface{}, error) {
	sfu.wg = wg
	sfu.ctx = ctx
	return &sfu.settings, nil
}

func (sfu *sfu) InitCommand() ([]*cobra.Command, error) {

	return nil, nil
}

func (sfu *sfu) ConfigChanged() {
}

func (sfu *sfu) RootCommand(cmd *cobra.Command, args []string) {
	sfu.wg.Add(1)

	sfu.ss = NewSignalServer(sfu.ctx, SignalServerSettings{
		listenAddr: sfu.settings.ListenAddr,
		cert:       sfu.settings.Cert,
		key:        sfu.settings.Key,
	})
	sfu.ss.OnClose(func() {
		sfu.wg.Done()
	})

	go func() {
		if e := sfu.ss.Run(); e != nil {
			logger().Errorf("signal server run failed: %v", e)
			panic(e)
		}
	}()
}

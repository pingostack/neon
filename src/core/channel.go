package core

import (
	"github.com/let-light/neon/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ChannelFlags struct {
}

type ChannelSettings struct {
}

type ChannelModule struct {
	flags    *ChannelFlags
	settings *ChannelSettings
}

var insChannel *ChannelModule

func init() {
	insChannel = &ChannelModule{
		flags:    &ChannelFlags{},
		settings: &ChannelSettings{},
	}
	module.Register(insChannel)
}

func ChannelModuleInstance() module.IModule {
	return insChannel
}

func (c *ChannelModule) OnInitModule() (interface{}, error) {
	return c.settings, nil
}

func (c *ChannelModule) OnInitCommand() ([]*cobra.Command, error) {
	return nil, nil
}

func (c *ChannelModule) OnConfigModified() {

}

func (c *ChannelModule) OnPostInitCommand() {

}

func (c *ChannelModule) OnMainRun(cmd *cobra.Command, args []string) {
	logrus.Info("channel main running ...")
}

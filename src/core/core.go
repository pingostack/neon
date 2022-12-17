package core

import (
	"github.com/let-light/neon/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type CoreFlags struct {
}

type CoreSettings struct {
}

type CoreModule struct {
	flags    *CoreFlags
	settings *CoreSettings
}

var instance *CoreModule

func init() {
	instance = &CoreModule{
		flags:    &CoreFlags{},
		settings: &CoreSettings{},
	}
	module.Register(instance)
}

func CoreModuleInstance() module.IModule {
	return instance
}

func (c *CoreModule) OnInitModule() (interface{}, error) {
	return c.settings, nil
}

func (c *CoreModule) OnInitCommand() ([]*cobra.Command, error) {
	return nil, nil
}

func (c *CoreModule) OnConfigModified() {

}

func (c *CoreModule) OnPostInitCommand() {

}

func (c *CoreModule) OnMainRun(cmd *cobra.Command, args []string) {
	logrus.Info("core main running ...")
}

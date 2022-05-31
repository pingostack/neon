package module

import (
	"github.com/spf13/cobra"
)

var cm *ConfigModule

type ConfigModule struct {
	output string
}

func init() {
	cm = &ConfigModule{}
}

func ConfigModuleInstance() *ConfigModule {
	return cm
}

func (c *ConfigModule) OnInitModule() (error, interface{}) {
	return nil, nil
}

func (c *ConfigModule) OnInitCommand() (error, []*cobra.Command) {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Operation configuration file",
		Long:    `Operation configuration file.`,
		Args:    cobra.MinimumNArgs(1),
		Example: `xxx config -f [config-file]`,
		Run:     c.run,
	}

	cmd.Flags().StringVarP(&c.output, "export", "e", "", "Export the configuration template file")

	return nil, []*cobra.Command{cmd}
}

func (c *ConfigModule) OnConfigModified() {
}

func (c *ConfigModule) OnPostInitCommand() {

}

func (c *ConfigModule) run(cmd *cobra.Command, args []string) {
}

func (c *ConfigModule) OnMainRun(cmd *cobra.Command, args []string) {
}

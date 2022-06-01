package module

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cm *ConfigModule

type ConfigModule struct {
	export string
	file   string
	consul string
	etcd   string
	config *viper.Viper
}

func init() {
	cm = &ConfigModule{}
}

func ConfigModuleInstance() *ConfigModule {
	return cm
}

func (c *ConfigModule) OnInitModule() (interface{}, error) {
	return nil, nil
}

func (c *ConfigModule) OnInitCommand() ([]*cobra.Command, error) {
	GetRootCmd().PersistentFlags().StringVarP(&c.file, "config", "c", "neon.yaml", "Load config file")
	GetRootCmd().PersistentFlags().StringVar(&c.consul, "consul", "", "Load config file from consul")
	GetRootCmd().PersistentFlags().StringVar(&c.etcd, "etcd", "", "Load config file from etcd")

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Operation configuration file",
		Long:  `Operation configuration file.`,
		Run:   c.run,
	}

	cmd.Flags().StringVarP(&c.export, "export", "e", "", "Export the configuration template file")

	return []*cobra.Command{cmd}, nil
}

func (c *ConfigModule) OnConfigModified() {
}

func (c *ConfigModule) OnPostInitCommand() {
	if c.consul != "" {
		c.loadConfigFromConsul()
	} else if c.etcd != "" {
		c.loadConfigFromEtcd()
	} else {
		c.loadConfigFromFile()
	}
}

func (c *ConfigModule) loadConfigFromFile() {
	c.config = viper.New()
	c.config.SetConfigFile(c.file)
	err := c.config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file, %s", err))
	}
}

func (c *ConfigModule) loadConfigFromConsul() {

}

func (c *ConfigModule) loadConfigFromEtcd() {

}

func (c *ConfigModule) run(cmd *cobra.Command, args []string) {
}

func (c *ConfigModule) OnMainRun(cmd *cobra.Command, args []string) {
}

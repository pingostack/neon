package module

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cm *ConfigModule

type ConfigFlags struct {
	ConfigFile string `env:"config" flag:"config"`
	Consul     string `env:"consul" flag:"consul"`
	Etcd       string `env:"etcd"   flag:"etcd"`
}

type ConfigModule struct {
	flags  ConfigFlags
	config *viper.Viper
}

func init() {
	cm = &ConfigModule{}
}

func ConfigModuleInstance() *ConfigModule {
	return cm
}

func (c *ConfigModule) Viper() *viper.Viper {
	return c.config
}

func (c *ConfigModule) OnInitModule() (interface{}, error) {
	return nil, nil
}

func (c *ConfigModule) OnInitCommand() ([]*cobra.Command, error) {
	GetRootCmd().PersistentFlags().StringVarP(&c.flags.ConfigFile, "config", "c", "config/neon.yaml", "Load config file")
	GetRootCmd().PersistentFlags().StringVar(&c.flags.Consul, "consul", "", "Load config file from consul")
	GetRootCmd().PersistentFlags().StringVar(&c.flags.Etcd, "etcd", "", "Load config file from etcd")

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Operation configuration file",
		Long:  `Operation configuration file.`,
		Run:   c.run,
	}

	return []*cobra.Command{cmd}, nil
}

func (c *ConfigModule) OnConfigModified() {
}

func (c *ConfigModule) OnPostInitCommand() {
	if c.flags.Consul != "" {
		c.loadConfigFromConsul()
	} else if c.flags.Etcd != "" {
		c.loadConfigFromEtcd()
	} else {
		c.loadConfigFromFile()
	}
}

func (c *ConfigModule) loadConfigFromFile() {
	c.config = viper.New()
	c.config.SetConfigFile(c.flags.ConfigFile)
	err := c.config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file, %s", err))
	}

	reloadSettings()
}

func (c *ConfigModule) loadConfigFromConsul() {

}

func (c *ConfigModule) loadConfigFromEtcd() {
}

// ./neon config reload ${config-file} -> reload config file
// ./neon config init ${output-config-file} -> init config file
func (c *ConfigModule) run(cmd *cobra.Command, args []string) {
}

func (c *ConfigModule) OnMainRun(cmd *cobra.Command, args []string) {
}

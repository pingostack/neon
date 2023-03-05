package gomodule

import (
	"context"
	"fmt"
	"sync"

	"github.com/gogf/gf/os/gfile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configInstance *configModule

type ConfigFlags struct {
	ConfigFile string `env:"config" flag:"config"`
	Consul     string `env:"consul" flag:"consul"`
	Etcd       string `env:"etcd"   flag:"etcd"`
}

type configModule struct {
	flags  ConfigFlags
	config *viper.Viper
	wg     *sync.WaitGroup
}

func init() {
	configInstance = &configModule{}
}

func ConfigModule() *configModule {
	return configInstance
}

func (c *configModule) Viper() *viper.Viper {
	return c.config
}

func (c *configModule) InitModule(ctx context.Context, wg *sync.WaitGroup) (interface{}, error) {
	c.wg = wg
	return nil, nil
}

func (c *configModule) InitCommand() ([]*cobra.Command, error) {
	GetRootCmd().PersistentFlags().StringVarP(&c.flags.ConfigFile, "config", "c", "config.yml", "Load config file")
	GetRootCmd().PersistentFlags().StringVar(&c.flags.Consul, "consul", "", "Load config file from consul")
	GetRootCmd().PersistentFlags().StringVar(&c.flags.Etcd, "etcd", "", "Load config file from etcd")

	return nil, nil
}

func (c *configModule) ConfigChanged() {
}

func (c *configModule) loadConfigFromFile() {
	if !gfile.Exists(c.flags.ConfigFile) {
		panic(fmt.Errorf("fatal error config file, %s not found", c.flags.ConfigFile))
	}

	c.config = viper.New()
	c.config.SetConfigFile(c.flags.ConfigFile)
	err := c.config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file, %s", err))
	}

	reloadSettings()
}

func (c *configModule) loadConfigFromConsul() {

}

func (c *configModule) loadConfigFromEtcd() {
}

func (c *configModule) RootCommand(cmd *cobra.Command, args []string) {
	if c.flags.Consul != "" {
		c.loadConfigFromConsul()
	} else if c.flags.Etcd != "" {
		c.loadConfigFromEtcd()
	} else {
		c.loadConfigFromFile()
	}
}

func reloadSettings() {
	for _, mi := range manager.modules {
		if mi == nil || mi.settings == nil {
			continue
		}

		if err := ConfigModule().Viper().UnmarshalKey(mi.name, mi.settings); err != nil {
			panic(fmt.Errorf("unmarshal config error, %s", err))
		}
	}
}

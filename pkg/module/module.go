package module

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/spf13/cobra"
)

var manager *Manager

type ModuleInfo struct {
	module   IModule
	settings interface{}
	cmds     []*cobra.Command
	name     string
}

type Manager struct {
	modules []*ModuleInfo
	rootCmd *cobra.Command
	once    sync.Once
}

type IModule interface {
	OnInitModule() (interface{}, error)
	OnInitCommand() ([]*cobra.Command, error)
	OnConfigModified()
	OnPostInitCommand()
	OnMainRun(cmd *cobra.Command, args []string)
}

func init() {
	manager = NewManager()
}

func NewManager() *Manager {
	m := &Manager{
		modules: make([]*ModuleInfo, 0),
		rootCmd: &cobra.Command{},
	}

	m.rootCmd.Run = func(cmd *cobra.Command, args []string) {
		for _, mi := range manager.modules {
			mi.module.OnMainRun(cmd, args)
		}
	}

	return m
}

func initDefaultModule() {
	Register(ConfigModuleInstance())
	Register(LoggerModuleInstance())
}

func Register(module IModule) error {
	if module == nil {
		return fmt.Errorf("module is nil")
	}

	t := reflect.TypeOf(module)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("module must be pointer")
	}

	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("module must be struct")
	}

	for _, mi := range manager.modules {
		if mi.module == module {
			return fmt.Errorf("module[%p] is existed", module)
		}
	}

	mi := &ModuleInfo{
		module: module,
		cmds:   make([]*cobra.Command, 0),
		name:   t.Elem().Name(),
	}

	cobra.OnInitialize(func() {
		module.OnPostInitCommand()
	})

	manager.modules = append(manager.modules, mi)

	return nil
}

func Launch() error {
	manager.once.Do(initDefaultModule)

	// init module
	for _, mi := range manager.modules {
		settings, err := mi.module.OnInitModule()
		if err != nil {
			return err
		}
		mi.settings = settings
	}

	// init command
	for _, mi := range manager.modules {
		cmds, err := mi.module.OnInitCommand()
		if err != nil {
			return err
		}

		for _, cmd := range cmds {
			manager.rootCmd.AddCommand(cmd)
		}

		mi.cmds = cmds
	}

	manager.rootCmd.Execute()

	return nil
}

func GetRootCmd() *cobra.Command {
	return manager.rootCmd
}

func reloadSettings() {
	for _, mi := range manager.modules {
		if mi == nil || mi.settings == nil {
			continue
		}

		if err := ConfigModuleInstance().Viper().UnmarshalKey(mi.name, mi.settings); err != nil {
			panic(fmt.Errorf("unmarshal config error, %s", err))
		}
	}
}

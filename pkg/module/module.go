package module

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"
)

var manager *Manager

type ModuleInfo struct {
	module IModule
	config interface{}
	cmds   []*cobra.Command
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

	m.rootCmd.Run = run

	return m
}

func initDefaultModule() {
	Register(ConfigModuleInstance())
}

func Register(module IModule) error {
	if module == nil {
		return fmt.Errorf("module is nil")
	}

	for _, mi := range manager.modules {
		if mi.module == module {
			return fmt.Errorf("module[%p] is existed", module)
		}
	}

	cobra.OnInitialize(module.OnPostInitCommand)
	mi := &ModuleInfo{
		module: module,
		cmds:   make([]*cobra.Command, 0),
	}

	manager.modules = append(manager.modules, mi)

	return nil
}

func Launch() error {
	manager.once.Do(initDefaultModule)

	// init module
	for _, mi := range manager.modules {
		config, err := mi.module.OnInitModule()
		if err != nil {
			return err
		}
		mi.config = config
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

func run(cmd *cobra.Command, args []string) {
	for _, mi := range manager.modules {
		mi.module.OnMainRun(cmd, args)
	}
}

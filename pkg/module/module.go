package module

import (
	"fmt"

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
}

type IModule interface {
	OnInitModule() (error, interface{})
	OnInitCommand() (error, []*cobra.Command)
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

func Init() {
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
	// init module
	for _, mi := range manager.modules {
		err, config := mi.module.OnInitModule()
		if err != nil {
			return err
		}
		mi.config = config
	}

	// init command
	for _, mi := range manager.modules {
		err, cmds := mi.module.OnInitCommand()
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

package gomodule

import (
	"context"
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
	modules        []*ModuleInfo
	rootCmd        *cobra.Command
	once           sync.Once
	ctx            context.Context
	cancel         context.CancelFunc
	wg             *sync.WaitGroup
	defaultModules []IModule
}

type IModule interface {
	InitModule(ctx context.Context, wg *sync.WaitGroup) (interface{}, error)
	InitCommand() ([]*cobra.Command, error)
	ConfigChanged()
	RootCommand(cmd *cobra.Command, args []string)
}

func init() {
	manager = NewManager()
}

func NewManager() *Manager {
	m := &Manager{
		modules:        make([]*ModuleInfo, 0),
		rootCmd:        &cobra.Command{},
		defaultModules: make([]IModule, 0),
	}

	m.wg = &sync.WaitGroup{}
	m.ctx, m.cancel = context.WithCancel(context.Background())

	m.rootCmd.Run = func(cmd *cobra.Command, args []string) {
		for _, mi := range manager.modules {
			mi.module.RootCommand(cmd, args)
		}
	}

	return m
}

func initDefaultModules() {
	modules := make([]*ModuleInfo, 0)
	for _, module := range manager.defaultModules {
		if module == nil {
			fmt.Printf("module is nil")
			continue
		}

		t := reflect.TypeOf(module)
		if t.Kind() != reflect.Ptr {
			fmt.Printf("module must be pointer")
			continue
		}

		if t.Elem().Kind() != reflect.Struct {
			fmt.Printf("module must be struct")
			continue
		}

		for _, mi := range manager.modules {
			if mi.module == module {
				fmt.Printf("module[%p] is existed", module)
				return
			}
		}

		mi := &ModuleInfo{
			module: module,
			cmds:   make([]*cobra.Command, 0),
			name:   t.Elem().Name(),
		}

		modules = append(modules, mi)
	}

	manager.modules = append(modules, manager.modules...)
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

	manager.modules = append(manager.modules, mi)

	return nil
}

func RegisterDefaultModule(modules ...IModule) {
	manager.defaultModules = append(manager.defaultModules, modules...)
}

func Launch(ctx context.Context) error {
	manager.ctx, manager.cancel = context.WithCancel(ctx)
	manager.once.Do(initDefaultModules)

	// init module
	for _, mi := range manager.modules {
		settings, err := mi.module.InitModule(ctx, manager.wg)
		if err != nil {
			return err
		}
		mi.settings = settings
	}

	// init command
	for _, mi := range manager.modules {
		cmds, err := mi.module.InitCommand()
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

func Wait() {
	go func() {
		manager.wg.Wait()
		manager.cancel()
	}()

	<-manager.ctx.Done()
}

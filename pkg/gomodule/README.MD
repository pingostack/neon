# example

```golang
package main

import (
	"context"
	"sync"

	"github.com/let-light/gomodule"
	"github.com/spf13/cobra"
)

type MainFlags struct {
}

type MainSettings struct {
}

type MainModule struct {
	flags    *MainFlags
	settings *MainSettings
	wg       *sync.WaitGroup
}

var instance *MainModule

func init() {
	instance = &MainModule{
		flags:    &MainFlags{},
		settings: &MainSettings{},
	}
	gomodule.Register(instance)
}

func MainModuleInstance() gomodule.IModule {
	return instance
}

func (c *MainModule) InitModule(ctx context.Context, wg *sync.WaitGroup) (interface{}, error) {
	c.wg = wg
	return c.settings, nil
}

func (c *MainModule) InitCommand() ([]*cobra.Command, error) {
	return nil, nil
}

func (c *MainModule) ConfigChanged() {

}

func (c *MainModule) CommandRun() {

}

func (c *MainModule) RootCommand(cmd *cobra.Command, args []string) {
}

func main() {
	gomodule.Launch(context.Background())
	gomodule.Wait()
}
```
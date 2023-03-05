package server

import (
	"context"
	"sync"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/let-light/neon/pkg/gomodule"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverInstance *serverModule

type serverModule struct {
	wg  *sync.WaitGroup
	ctx context.Context
}

func init() {
	serverInstance = &serverModule{}
	gomodule.Register(serverInstance)
}

func ServerModule() *serverModule {
	return serverInstance
}

func (server *serverModule) InitModule(ctx context.Context, wg *sync.WaitGroup) (interface{}, error) {
	server.wg = wg
	server.ctx = ctx
	return nil, nil
}

func (server *serverModule) InitCommand() ([]*cobra.Command, error) {

	return nil, nil
}

func (server *serverModule) ConfigChanged() {
}

func (server *serverModule) RootCommand(cmd *cobra.Command, args []string) {
	server.wg.Add(1)
	addr := ":8888"
	s1 := g.Server("test")
	s1.SetAddr(addr)
	group := s1.Group("v1")
	group.ALL("test", func(r *ghttp.Request) {
		logrus.Info("test")
		r.Response.Write("test")
	})

	s1.Start()

	go func() {
		<-server.ctx.Done()
		logrus.Info("server shutdown")
		s1.Shutdown()
	}()

	go func() {
		logrus.Info("server start at ", addr)
		g.Wait()
		server.wg.Done()
		logrus.Info("server stop")
	}()
}

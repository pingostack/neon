package server

import (
	"context"
	"sync"

	"github.com/gogf/gf/os/gfile"
	"github.com/kardianos/service"
	"github.com/let-light/neon/pkg/gomodule"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type syservice struct {
	ctx    context.Context
	cancel context.CancelFunc
	svc    service.Service
}

var serviceInstance *syservice

func init() {

	ctx, cancel := context.WithCancel(context.Background())
	serviceInstance = &syservice{
		ctx:    ctx,
		cancel: cancel,
	}
}

func ServiceModule() *syservice {
	return serviceInstance
}

func (s *syservice) Start(ss service.Service) error {
	return nil
}

func (s *syservice) Stop(ss service.Service) error {
	logrus.Info("service stop")
	s.cancel()
	return nil
}

func (s *syservice) InitModule(ctx context.Context, wg *sync.WaitGroup) (interface{}, error) {

	workdir := gfile.Join(gfile.Dir(gfile.SelfDir()))
	gfile.Chdir(workdir)
	configFile := gfile.Join(workdir, "config", "config.yml")
	logrus.Infof("workdir: %s", workdir)

	var err error
	s.svc, err = service.New(s, &service.Config{
		Name:             "neon",
		DisplayName:      "neon",
		Description:      "neon service",
		WorkingDirectory: workdir,
		Arguments:        []string{"--config", configFile},
	})

	if err != nil {
		logrus.Error("new service failed:", err)
	}

	return nil, err
}

func (s *syservice) InitCommand() ([]*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "service",
		Short: `[install|uninstall|start|stop]`,
		Long:  `it is a command for service install uninstall start stop`,
		Run: func(cmd *cobra.Command, args []string) {
			install, err := cmd.Flags().GetBool("install")
			if err != nil {
				logrus.Error("service command get install flag error:", err)
				return
			}

			uninstall, err := cmd.Flags().GetBool("uninstall")
			if err != nil {
				logrus.Error("service command get uninstall flag error:", err)
				return
			}

			start, err := cmd.Flags().GetBool("start")
			if err != nil {
				logrus.Error("service command get start flag error:", err)
				return
			}

			stop, err := cmd.Flags().GetBool("stop")
			if err != nil {
				logrus.Error("service command get stop flag error:", err)
				return
			}

			restart, err := cmd.Flags().GetBool("restart")
			if err != nil {
				logrus.Error("service command get restart flag error:", err)
				return
			}

			if install {
				err = s.svc.Install()
				if err != nil {
					logrus.Errorf("service install error: ", err)
					return
				}
				logrus.Infof("service install success")
			}

			if uninstall {
				err = s.svc.Uninstall()
				if err != nil {
					logrus.Errorf("service uninstall error: ", err)
					return
				}
				logrus.Infof("service uninstall success")
				return
			}

			if start {
				err = s.svc.Start()
				if err != nil {
					logrus.Errorf("service start error: ", err)
					return
				}
				logrus.Infof("service start success")
			}

			if stop {
				err = s.svc.Stop()
				if err != nil {
					logrus.Errorf("service stop error: ", err)
					return
				}
				logrus.Infof("service stop success")
			}

			if restart {
				err = s.svc.Restart()
				if err != nil {
					logrus.Errorf("service restart error: ", err)
					return
				}
				logrus.Infof("service restart success")
			}
		},
	}

	cmd.Flags().BoolP("install", "i", false, "install neon service")
	cmd.Flags().BoolP("uninstall", "u", false, "uninstall neon service")
	cmd.Flags().BoolP("start", "s", false, "start neon service")
	cmd.Flags().BoolP("stop", "t", false, "stop neon service")
	cmd.Flags().BoolP("restart", "r", false, "restart neon service")

	return []*cobra.Command{cmd}, nil

}

func (s *syservice) ConfigChanged() {

}

func (s *syservice) RootCommand(cmd *cobra.Command, args []string) {
	go s.svc.Run()
}

func Launch() {
	gomodule.RegisterDefaultModule(gomodule.ConfigModule())
	gomodule.RegisterDefaultModule(gomodule.LoggerModule())
	gomodule.RegisterDefaultModule(ServiceModule())
	gomodule.Register(ServerModule())

	gomodule.Launch(ServiceModule().ctx)
}

func Wait() {
	gomodule.Wait()
	logrus.Info("neon server exit")
}

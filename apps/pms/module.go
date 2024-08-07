package pms

import (
	"context"
	"sync"
	"time"

	"github.com/let-light/gomodule"
	pms_feature "github.com/pingostack/neon/features/pms"
	"github.com/pingostack/neon/internal/httpserv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pmsModule *pms

type ISignalServer interface {
	Start() error
	Close() error
}

type PMSSettings struct {
	httpserv.HttpParams    `json:"http" mapstructure:"http"`
	KeyFrameIntervalSecond time.Duration `json:"keyFrameIntervalSeconds" mapstructure:"keyFrameIntervalSeconds"`
	JoinTimeoutSecond      time.Duration `json:"joinTimeoutSeconds" mapstructure:"joinTimeoutSeconds"`
}

type pms struct {
	gomodule.DefaultModule
	ctx         context.Context
	preSettings PMSSettings
	settings    *PMSSettings
	logger      *logrus.Entry
	serv        ISignalServer
	lock        sync.RWMutex
}

func init() {
	pmsModule = &pms{
		logger: logrus.WithField("module", "pms"),
	}
}

func PMSModule() gomodule.IModule {
	return pmsModule
}

func settings() PMSSettings {
	return pmsModule.Settings()
}

func (pms *pms) InitModule(ctx context.Context, _ *gomodule.Manager) (interface{}, error) {
	pms.ctx = ctx
	return &pms.preSettings, nil
}

func (pms *pms) InitCommand() ([]*cobra.Command, error) {

	return nil, nil
}

func (pms *pms) ConfigChanged() {
	if pms.settings == nil {
		pms.lock.Lock()
		defer pms.lock.Unlock()
		pms.settings = &pms.preSettings
	}
}

func (pms *pms) Settings() PMSSettings {
	pms.lock.RLock()
	defer pms.lock.RUnlock()
	return *pms.settings
}

func (pms *pms) ModuleRun() {
	pms.serv = NewSignalServer(pms.ctx, pms.logger)
	if err := pms.serv.Start(); err != nil {
		pms.logger.Errorf("pms start error: %v", err)
		panic(errors.Wrap(err, "pms start error"))
	}

	<-pms.ctx.Done()
	pms.close()
}

func (pms *pms) Type() interface{} {
	return pms_feature.Type()
}

func (pms *pms) close() {
	pms.logger.Info("pms closing")
	pms.serv.Close()
}

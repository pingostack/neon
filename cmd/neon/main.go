package main

import (
	"runtime"
	"time"

	"github.com/kardianos/service"
	"github.com/let-light/neon/cmd/neon/server"
	"github.com/sirupsen/logrus"
)

var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() error {
	logger.Infof("I'm running %v.", service.Platform())
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case tm := <-ticker.C:
			logger.Infof("Still running at %v...", tm)
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
}
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

func PrintStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], true)
	logrus.Errorf("==> %s\n", string(buf[:n]))
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Error during:", err)
			PrintStack()
		}
	}()

	server.Launch()
	server.Wait()

}

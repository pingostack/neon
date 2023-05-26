package main

import (
	"context"
	"runtime"

	"github.com/let-light/gomodule"
	"github.com/let-light/neon/apps/sfu"
	"github.com/let-light/neon/pkg/gortc"
	"github.com/sirupsen/logrus"
)

func printStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], true)
	logrus.Errorf("==> %s\n", string(buf[:n]))
}

func serv(ctx context.Context) {
	gomodule.RegisterDefaultModules()
	gomodule.RegisterWithName(gortc.WebRTCModule(), "webrtc")
	gomodule.RegisterWithName(sfu.SfuModule(), "sfu")

	gomodule.Launch(ctx)

	gomodule.Wait()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("Error during:", err)
			printStack()
		}
	}()

	serv(context.Background())
}

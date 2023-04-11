package main

import (
	"runtime"

	"github.com/let-light/neon/cmd/neon/server"
	"github.com/sirupsen/logrus"
)

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

/*
Copyright © 2022 im-pingo HERE <cczjp89@gmail>

*/
package main

import (
	"fmt"
	"time"

	"github.com/let-light/neon/pkg/module"
	_ "github.com/let-light/neon/src/core"

	_ "github.com/let-light/neon/src/modules/rtsp"
	_ "github.com/let-light/neon/src/modules/webrtc"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	module.Launch()

	for {
		time.Sleep(time.Second * 1)
	}
}

/*
Copyright © 2022 im-pingo HERE <cczjp89@gmail>

*/
package main

import (
	"fmt"

	"github.com/pingopenstack/neon/pkg/module"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	module.Launch()

	logrus.Info("neon exit")
}

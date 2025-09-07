package main

import (
	"github.com/spf13/viper"
	"github.com/umutcomlekci/automated-messaging-system/cmd"
)

func main() {
	viper.AutomaticEnv()
	if err := cmd.NewCommand().Execute(); err != nil {
		panic(err)
	}
}

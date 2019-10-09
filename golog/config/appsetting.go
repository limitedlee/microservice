package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

var App AppSetting

func init() {
	_, err := toml.DecodeFile("appsetting.toml", &App)
	if err != nil {
		log.Fatal(err)
	}
}

type AppSetting struct {
	Grpc struct {
		Address string
		Appid   string
	}
}

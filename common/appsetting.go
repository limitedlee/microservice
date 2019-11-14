package common

import (
	"github.com/BurntSushi/toml"
	"log"
)

var PbConfig SysConfig

type SysConfig struct {
	Grpc struct {
		Appid   string //項目id
		Address string //配置服務的grpc地址
	}
}

func init() {
	_, err := toml.DecodeFile("appsetting.toml", &PbConfig)
	if err != nil {
		log.Fatal(err)
	}
}

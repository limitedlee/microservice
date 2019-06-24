package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/limitedlee/microservice/common"
	"github.com/limitedlee/microservice/consul"
	"github.com/limitedlee/microservice/jwt"
	"github.com/limitedlee/microservice/rsa"
	"strings"
)

func main() {

	sysConfig := &common.SystemConfig{
		Name:                    "godemo",
		DisplayName:             "go语言的demo",
		LocalAddress:            "10.1.1.51:9527",
		ServiceDiscoveryAddress: "test.dc.rpdns.com:8500"}

	consul.RegisterService(sysConfig)

	publickey, _ := rsa.LoadRsaKey("rsa_1024_pub.pem", "rsa_1024_priv.pem")
	jwt.PublicKey = publickey

	route := gin.Default()
	route.Use(jwt.JWT())
	route.GET("/health", healthCheck)

	port := strings.Split(sysConfig.LocalAddress, ":")[1]

	_ = route.Run(fmt.Sprintf(":%s", port))
}

func healthCheck(ctx *gin.Context) {
	ctx.String(200, "ok")
}

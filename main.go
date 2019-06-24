package microservice

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/LimitedLee/microservice/common"
	consulHelper "github.com/LimitedLee/microservice/consul"
	"github.com/LimitedLee/microservice/jwt"
	"github.com/LimitedLee/microservice/rsa"
	"strings"
)

func main() {

	sysConfig := &common.SystemConfig{
		Name:                    "godemo",
		DisplayName:             "go语言的demo",
		LocalAddress:            "10.1.1.51:9527",
		ServiceDiscoveryAddress: "test.dc.rpdns.com:8500"}

	consulHelper.RegisterService(sysConfig)

	publickey,_:=rsa.LoadRsaKey("rsa_1024_pub.pem","rsa_1024_priv.pem")
	jwt.PublicKey=publickey

	route := gin.Default()
	route.Use(jwt.JWT())
	route.GET("/health", healthCheck)
	route.GET("/login",login)

	port := strings.Split(sysConfig.LocalAddress, ":")[1]

	route.Run(fmt.Sprintf(":%s",port))
}

func healthCheck(ctx *gin.Context) {
	ctx.String(200, "ok")
}

func login(ctx *gin.Context)  {

}

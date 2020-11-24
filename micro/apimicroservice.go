package micro

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/limitedlee/microservice/common/config"
	"github.com/limitedlee/microservice/common/nacos"
	"github.com/lsls907/nacos-sdk-go/vo"
	"strconv"
	"strings"
)

type ApiMicroService struct {
	echo.Echo
}

func (a *ApiMicroService) NewServer(){
	 echo.New()
}

func (a *ApiMicroService) StartApi(serviceName string) error {
	baseUrl, _ := config.Get("BaseUrl")
	items := strings.Split(baseUrl, ":")
	addr := fmt.Sprintf(":%v", items[len(items)-1])

	if len(items) <= 0 {
		panic("Please define the port，example(:7065)")
	}
	intNum, _ := strconv.Atoi(items[1])
	port := uint64(intNum)

	nacos.RegisterServiceInstance(vo.RegisterInstanceParam{
		Ip:          nacos.GetOutboundIp(),
		Port:        port,
		ServiceName: serviceName,
		Weight:      10,
		ClusterName: "DEFAULT",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		GroupName:   "DEFAULT_GROUP",
	})
	return a.Start(addr)
}

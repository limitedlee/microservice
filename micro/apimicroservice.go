package micro

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/limitedlee/microservice/common/config"
	"github.com/limitedlee/microservice/common/handles"
	"github.com/limitedlee/microservice/common/nacos"
	"github.com/lsls907/nacos-sdk-go/vo"
	"strconv"
	"strings"
)

type ApiMicroService struct {
	//echo.Echo
}

func (a *ApiMicroService) NewServer() *echo.Echo {
	return echo.New()
}

//注入nacos
func (a *ApiMicroService) StartApi(e *echo.Echo, serviceName string, addr string) error {
	r := e.Group("/pool")
	r.POST("/change", handles.ApiChangesPool)
	port, addr := getAddr(addr)

	nacos.RegisterServiceInstance(vo.RegisterInstanceParam{
		Ip:          nacos.GetOutboundIp(),
		Port:        port,
		ServiceName: serviceName,
		Weight:      10,
		ClusterName: "DEFAULT",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		GroupName:   "HTTP",
	})
	return e.Start(addr)
}

func getAddr(addr string) (uint64, string) {
	baseUrl := ""
	if len(addr) <= 0 {
		baseUrl, _ = config.Get("BaseUrl")
	} else {
		baseUrl = addr
	}
	items := strings.Split(baseUrl, ":")
	if len(items) <= 0 {
		panic("Please define the port，example(:7065)")
	}
	port, _ := strconv.Atoi(items[1])
	addr = fmt.Sprintf(":%v", items[len(items)-1])
	return uint64(port), addr

}

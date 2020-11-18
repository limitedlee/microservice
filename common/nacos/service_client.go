package nacos

import (
	"fmt"
	"github.com/limitedlee/microservice/common/logger"
	"github.com/limitedlee/microservice/common/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
	"strconv"
)

//Register with default cluster and group
//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
func RegisterServiceInstance(param vo.RegisterInstanceParam) {
	ipAddr, _ := config.Get("nacos-addr")
	port, _ := config.Get("nacos-port")
	namespaceId, _ := config.Get("nacos-namespace-id")
	intPort, _ := strconv.Atoi(port)
	sc := []constant.ServerConfig{
		{
			IpAddr: ipAddr,
			Port:   uint64(intPort),
		},
	}
	cc := constant.ClientConfig{
		NamespaceId:         namespaceId, //namespace id
		TimeoutMs:           5000,        // 请求Nacos服务端的超时时间，默认是10000ms
		NotLoadCacheAtStart: true,        // 在启动的时候不读取缓存在CacheDir的service信息
		//LogDir:              "/tmp/nacos/log", // 日志存储路径
		//CacheDir:            "/tmp/nacos/cache",
		RotateTime: "24h",  // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
		MaxAge:     3,      // 日志最大文件数，默认3
		LogLevel:   "info", // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
	}

	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)

	}
	success, _ := client.RegisterInstance(param)
	logger.Info(fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success))
}

// Get preferred outbound ip of this machine
func GetOutboundIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

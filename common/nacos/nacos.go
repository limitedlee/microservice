package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strconv"
	"strings"
	"sync"
	"time"
)

var client naming_client.INamingClient               //nacos客户端
var rwLock *sync.RWMutex                             //读写锁
var ServiceRoute map[string][]model.SubscribeService //本地服务缓存

//初始化服务发现
func InitDiscovery(dsAddress, namespaceId string) {
	address := strings.Split(dsAddress, ":")[0]
	port, _ := strconv.Atoi(strings.Split(dsAddress, ":")[1])

	sc := []constant.ServerConfig{
		{
			IpAddr: address,
			Port:   uint64(port),
			Scheme: "http",
		},
	}

	cc := constant.ClientConfig{
		TimeoutMs:   500,
		NamespaceId: namespaceId,
		//CacheDir:             "e:/nacos/cache",
		NotLoadCacheAtStart:  true,
		UpdateCacheWhenEmpty: false,
		//LogDir:               "e:/nacos/log",
	}

	client, _ = clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
}

//订阅
func SubServices() {
	//for{
	fmt.Println("-----------------------------------")
	//从nacos上拉取注册的服务列表
	services, _ := client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		PageNo:   1,
		PageSize: 9999,
	})

	//对列表中的每一个服务都订阅
	for _, v := range services.Doms {
		if _, ok := ServiceRoute[v]; !ok {
			client.Subscribe(&vo.SubscribeParam{
				ServiceName:       v,
				SubscribeCallback: changeNotity,
			})

			ServiceRoute[v] = nil

			fmt.Println("已订阅：", v)
		}
	}

	//	time.Sleep(10*time.Second)
	//}
}

//监控服务变化，及时订阅最新上线的服务
func WatchServices() {
	for true {
		time.Sleep(10 * time.Second)
	}
}

func changeNotity(subInstances []model.SubscribeService, err error) {
	//不处理空推送
	if len(subInstances) == 0 {
		return
	}

	fmt.Println("收到订阅:", len(subInstances))

	//遍历订阅服务列表
	for _, subInstance := range subInstances {
		////判断当前服务是否存在本地服务路由表中，如果不在则加入进去
		//if _, ok := ServiceRoute[subInstance.ServiceName]; !ok {
		//	ServiceRoute[subInstance.ServiceName] = append(ServiceRoute[subInstance.ServiceName], subInstance)
		//	continue
		//}

		//当前 根据服务名、IP、端口 确定该服务是否已登记，后续可考虑加入版本号
		hasThisInstance := false
		currentInstanceIndex := -1
		for i, currentLocalInstance := range ServiceRoute[subInstance.ServiceName] {
			if currentLocalInstance.ServiceName == subInstance.ServiceName && currentLocalInstance.Ip == subInstance.Ip && currentLocalInstance.Port == subInstance.Port {
				hasThisInstance = true
				currentInstanceIndex = i
				break
			}
		}
		//如果当前实例不存在就登记进来
		if !hasThisInstance && subInstance.Valid == true {
			ServiceRoute[subInstance.ServiceName] = append(ServiceRoute[subInstance.ServiceName], subInstance)
		}
		//如果已存在，且推送的状态为Valid，则移除掉
		if hasThisInstance && subInstance.Valid == false {
			beforeList := ServiceRoute[subInstance.ServiceName][:currentInstanceIndex]
			afterList := ServiceRoute[subInstance.ServiceName][currentInstanceIndex+1:]
			rwLock.Lock()
			ServiceRoute[subInstance.ServiceName] = append(beforeList, afterList...)
			rwLock.Unlock()
		}
	}
}

package nacos

import (
	"encoding/json"
	"fmt"
	"github.com/limitedlee/microservice/common/config"
	"github.com/limitedlee/microservice/common/handles"
	wr "github.com/mroth/weightedrand"
	"google.golang.org/grpc"
	"math/rand"
	"reflect"
	"time"
)

var clientMap = make(map[string]ClientConnection, 0)

type ClientConnection struct {
	Conn    *grpc.ClientConn
	Message string
}

func (a ClientConnection) isEmpty() bool {
	return reflect.DeepEqual(a, ClientConnection{})
}

func GetGrpcConn(serverName string) ClientConnection {
	clientConnection := ClientConnection{}
	url := getPoolUrl(serverName)
	if len(url) <= 0 {
		clientConnection.Message = fmt.Sprintf("%s,没有注册到nacos上，或者服务不存在,", serverName)
		//errStr :=
		return clientConnection
	}
	if len(clientMap) <= 0 || clientMap[url].isEmpty() {
		var client *grpc.ClientConn
		client, err := grpc.Dial(url, grpc.WithInsecure())
		if err != nil {
			clientConnection.Message = err.Error()
			return clientConnection
		}
		clientConnection.Conn = client
		clientConnection.Message = "success"

		clientMap[url] = clientConnection
		return clientConnection

	}
	return clientMap[url]
}

func getPoolUrl(serverName string) string {
	if len(handles.GrpcPool) <= 0 || len(handles.GrpcPool[serverName]) <= 0 {
		//handles.GrpcPool[serverName]=
		namespaceId, _ := config.Get("nacos-namespace-id")
		data := getConfigs(ConfigRequest{
			Tenant: namespaceId,
			DataId: serverName,
			Group:  "DEFAULT_GROUP",
		})
		if len(data) <= 0 {
			return ""
		}
		poolMap := make(map[string][]PoolUrl, 0)
		_ = json.Unmarshal([]byte(data), poolMap)

		if len(poolMap[serverName]) <= 0 {
			return ""
		}
		handles.Mutex.Lock() //对共享变量操作之前先加锁
		handles.GrpcPool[serverName] = poolMap[serverName]
		handles.Mutex.Unlock() //对共享变量操作完毕在解锁，
		return getAdUrl(poolMap[serverName])
	}
	return getAdUrl(handles.GrpcPool[serverName])

}

func getAdUrl(data []PoolUrl) string {
	choices := make([]wr.Choice, 0)
	for _, v := range data {
		if v.Weight == 0 {
			continue
		}
		choice := wr.Choice{
			Item:   v.Url,
			Weight: uint(v.Weight),
		}
		choices = append(choices, choice)
	}
	if len(choices) <= 0 {
		return ""
	}
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!
	chooser, _ := wr.NewChooser(choices...)
	result := chooser.Pick().(string)
	return result
}

package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/limitedlee/microservice/model"
	"github.com/limitedlee/microservice/proto"
	"google.golang.org/grpc"
)

var client proto.AppconfigClient

func init() {
	conn, err := grpc.Dial(model.PbConfig.Grpc.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	client = proto.NewAppconfigClient(conn)
}

//根据key获取配置信息
func Get(key string) (string, error) {
	envName := os.Getenv("ASPNETCORE_ENVIRONMENT")
	//传递的key格式为AppId:EnId:KV
	key = fmt.Sprintf("%s:%s:%s", envName, model.PbConfig.Grpc.Appid, key)
	log.Println(key)

	// 调用gRPC接口
	var param proto.Params
	param.Keys = key
	tr, err := client.GetAppConfig(context.Background(), &param)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("服务端响应: %s", tr.Data)
	if err != nil {
		log.Fatalf("get config info Unmarshal fail : %v", err)
	}
	return tr.Data, err
}

func GetString(key string) string {
	str, _ := Get(key)
	return str
}

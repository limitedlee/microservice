package config

import (
	"context"
	"fmt"
	"github.com/limitedlee/microservice/common"
	"google.golang.org/grpc"
	"log"
	"os"
)

var client AppconfigClient

func init() {
	conn, err := grpc.Dial(common.PbConfig.Grpc.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	client = NewAppconfigClient(conn)
}

//根据key获取配置信息
func Get(key string) (string, error) {
	envName := os.Getenv("ASPNETCORE_ENVIRONMENT")
	fmt.Println("环境变量值", envName)
	//传递的key格式为AppId:EnId:KV
	key = fmt.Sprintf("%s:%s:%s", envName, common.PbConfig.Grpc.Appid, key)
	log.Println(key)

	// 调用gRPC接口
	var param Params
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
	str, err := Get(key)
	if err != nil {
		log.Println("读取配置出错 ", err)
	}
	return str
}

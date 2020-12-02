package grpc

import (
	"context"
	"fmt"
	"github.com/limitedlee/microservice/common/config"
	"google.golang.org/grpc"
)

type BaseClient struct {
	Conn *grpc.ClientConn
	//Ctx  context.Context
	Cf context.CancelFunc
}

func GetBaseClient(serverName string) *BaseClient {
	var client BaseClient
	var err error
	key := fmt.Sprintf("grpc.%s", serverName)
	url := config.GetString(key)

	if client.Conn, err = grpc.Dial(url, grpc.WithInsecure()); err != nil {
		return nil
	} else {
		//client.RPC = ac.NewApmClientServiceClient(client.Conn)
		//client.Ctx = context.Background()
		//client.Ctx, client.Cf = context.WithTimeout(client.Ctx, time.Second*30)
		return &client
	}
}

func (a *BaseClient) Close() {
	a.Conn.Close()
	a.Cf()
}

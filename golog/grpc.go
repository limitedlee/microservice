package golog

import (
	"context"
	"log"
	"time"

	"github.com/limitedlee/microservice/golog/config"
	"github.com/limitedlee/microservice/golog/proto"
	"google.golang.org/grpc"
)

var client proto.LogClient
var ctx context.Context
var cancel context.CancelFunc

func init() {
	conn, err := grpc.Dial(config.App.Grpc.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	client = proto.NewLogClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
}

func Error(msg ...interface{}) {
	if len(msg) == 0 {
		return
	}
	m := proto.LogRequest{}
	for _, v := range msg {
		switch v.(type) {
		case string:
			m.Message = v.(string)
		case error:
			m.Exception = v.(error).Error()
		}
	}

	r, err2 := client.Error(ctx, &m)
	if err2 != nil {
		log.Fatalf("could not greet: %v", err2)
	}
	log.Println(r)
}

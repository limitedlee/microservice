package main

import (
	"github.com/limitedlee/microservice/common/handles"
	pb "github.com/limitedlee/microservice/example/proto"
	"github.com/limitedlee/microservice/micro"
)

func main() {
	micro := &micro.MicService{}
	micro.Routes["/ws"] = handles.WebSocketHandler
	micro.NewServer()
	pb.RegisterUserServiceServer(micro.GrpcServer, &pb.UserService{})
	micro.Start()
}

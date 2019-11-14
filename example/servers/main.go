package servers

import (
	pb "github.com/limitedlee/microservice/example/proto"
	"github.com/limitedlee/microservice/micro"
)

func main() {
	micro := &micro.MicService{}
	micro.NewServer()
	pb.RegisterUserServiceServer(micro.GrpcServer, &UserService{})
	micro.Start()
}

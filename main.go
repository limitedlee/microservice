package main

func main() {
	server := MicService{}
	server.NewServer()
	//pb.RegisterUserServiceServer(server.GrpcServer,&services.UserServiceServer{})
	server.Start()
}

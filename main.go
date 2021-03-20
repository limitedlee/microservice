package main

import (
	"github.com/labstack/echo/v4"
	"github.com/limitedlee/microservice/example/proto"
	"github.com/limitedlee/microservice/micro"
	"google.golang.org/grpc"
)

func main() {
	service := micro.NewService(
		micro.SetServiceType(micro.GRPC),
		micro.SetServiceName("aaaa"),
		micro.SetRunMode(micro.ByHost),
	)

	//m:= service.Group("/m")
	//m.Get("/a",a)
	//m.POST("/b",b)
	//service.Get("/k",k)

	service.Init()
	service.RegisterService("mse-e52dbdd6-p.nacos-ans.mse.aliyuncs.com:8848", "27fdefc2-ae39-41fd-bac4-9256acbf97bc")
	proto.RegisterUserServiceServer(service.Instance.(*grpc.Server), &proto.UserService{})
	service.Run()
}

func a(c echo.Context) error {
	return c.String(200, "a")
}

func b(c echo.Context) error {
	return c.String(200, "a")
}

func k(c echo.Context) error {
	return c.String(200, "a")
}

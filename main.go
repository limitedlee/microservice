package main

import "github.com/limitedlee/microservice/service"

func main() {
	s := service.Microservice{}
	route := s.InitService()
	route.Group("")
	s.Run()
}

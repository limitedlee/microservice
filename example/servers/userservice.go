package servers

import (
	"context"
	pb "github.com/limitedlee/microservice/example/proto"
)

type UserService struct {
}

//获取指定用户接口
func (s *UserService) Get(context.Context, *pb.SearchUserParam) (*pb.User, error) {
	u := &pb.User{
		Id:     1,
		Name:   "张三",
		Age:    18,
		Sex:    false,
		Mobile: "13800168000",
		IdCard: "420881198602220112"}
	return u, nil
}

//获取用户列表，分页
func (s *UserService) List(context.Context, *pb.SearchUserParam) (*pb.UserList, error) {
	list := &pb.UserList{
		Index:    0,
		PageSize: 1,
		Count:    1}

	list.Users[0] = &pb.User{
		Id:     1,
		Name:   "张三",
		Age:    18,
		Sex:    false,
		Mobile: "13800168000",
		IdCard: "420881198602220112"}

	return list, nil
}

//添加用户
func (s *UserService) Add(context.Context, *pb.UserAddInfo) (*pb.BaseResponse, error) {
	return &pb.BaseResponse{Code: 0, Message: "1111"}, nil
}

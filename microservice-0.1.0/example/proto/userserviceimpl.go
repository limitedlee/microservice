package proto

import (
	"context"
)

type UserService struct {
}

//获取指定用户接口
func (s *UserService) Get(context.Context, *SearchUserParam) (*User, error) {
	u := &User{
		Id:     1,
		Name:   "张三",
		Age:    18,
		Sex:    false,
		Mobile: "13800168000",
		IdCard: "420881198602220112"}
	return u, nil
}

//获取用户列表，分页
func (s *UserService) List(context.Context, *SearchUserParam) (*UserList, error) {
	list := &UserList{
		Index:    0,
		PageSize: 1,
		Count:    1}

	list.Users[0] = &User{
		Id:     1,
		Name:   "张三",
		Age:    18,
		Sex:    false,
		Mobile: "13800168000",
		IdCard: "420881198602220112"}

	return list, nil
}

//添加用户
func (s *UserService) Add(context.Context, *UserAddInfo) (*BaseResponse, error) {
	return &BaseResponse{Code: 0, Message: "1111"}, nil
}

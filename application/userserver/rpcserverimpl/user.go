package rpcserverimpl

import (
	"context"
	"errors"

	"micro-me/application/userserver/model"
	userPb "micro-me/application/userserver/protos"
)

type (
	UserRpcServer struct {
		userModel *model.MembersModel
	}
)

var (
	ErrNotFound = errors.New("用户不存在 ")
)


func NewUserRpcServer(userModel *model.MembersModel)*UserRpcServer{
	return &UserRpcServer{
		userModel: userModel,
	}
}

func (s *UserRpcServer) FindByToken(ctx context.Context, req *userPb.FindByTokenRequest, rsp *userPb.UserResponse) error {
	member, err := s.userModel.FindByToken(req.Token)
	if err != nil {
		return ErrNotFound
	}
	rsp.Token = member.Token
	rsp.Id = member.Id
	rsp.Password = member.Password
	rsp.Username = member.Username
	return nil
}

func (s *UserRpcServer) FindById(ctx context.Context, req *userPb.FindByIdRequest, rsp *userPb.UserResponse) error {
	member, err := s.userModel.FindById(req.Id)
	if err != nil {
		return ErrNotFound
	}
	rsp.Token = member.Token
	rsp.Id = member.Id
	rsp.Password = member.Password
	rsp.Username = member.Username
	return nil
}

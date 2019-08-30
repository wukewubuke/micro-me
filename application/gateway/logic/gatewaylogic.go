package logic

import (
	"context"
	"micro-me/application/common/baseerror"
	"micro-me/application/common/config"
	"micro-me/application/gateway/model"
	imPb "micro-me/application/imserver/protos"
	userPb "micro-me/application/userserver/protos"
	"time"
)

type (
	GatewayLogic struct {
		userRpcModel  userPb.UserService
		imRpcModel    imPb.ImService
		gatewayModel  *model.GateWayModel
		imAddressList []*config.ImRpcServer
	}

	SendRequest struct {
		FromToken string    `json:"fromToken" binding:"required"`
		ToToken   string    `json:"toToken" binding:"required"`
		Body      string    `json:"body" binding:"required"`
		Timestamp time.Time `json:"timestamp"`
	}
	SendResponse struct {
	}

	GetImAddressRequest struct {
		Token string `json:"token" binding:"required"`
	}
	GetImAddressResponse struct {
		Address string `json:"address"`
	}
)

var (
	SendMessageErr    = baseerror.NewBaseError("发送消息失败")
	UserNotFoundErr   = baseerror.NewBaseError("用户不存在")
	ImAddressErr      = baseerror.NewBaseError("请配置IM服务地址")
	AddGatewayErr     = baseerror.NewBaseError("添加gateway数据库失败")
	PublishMessageErr = baseerror.NewBaseError("发送消息到MQ失败")
	ImRpcModelMapErr  = baseerror.NewBaseError("没有此IM服务")
)

func NewGatewayLogic(userRpcModel userPb.UserService, gatewayModel *model.GateWayModel, imAddressList []*config.ImRpcServer, imRpcModel imPb.ImService) *GatewayLogic {
	return &GatewayLogic{
		userRpcModel:  userRpcModel,
		gatewayModel:  gatewayModel,
		imAddressList: imAddressList,
		imRpcModel:    imRpcModel,
	}
}

func (l *GatewayLogic) Send(r *SendRequest) (*SendResponse, error) {

	if _, err := l.userRpcModel.FindByToken(context.TODO(), &userPb.FindByTokenRequest{Token: r.ToToken}); err != nil {
		return nil, UserNotFoundErr
	}

	gateway, err := l.gatewayModel.FindByToken(r.ToToken)
	if err != nil {
		return nil, SendMessageErr
	}

	if gateway.Id <= 0 {
		return nil, SendMessageErr
	}

	//发送消息逻辑

	if _, err := l.imRpcModel.PublishMessage(context.TODO(), &imPb.PublishMessageRequest{
		FromToken: r.FromToken,
		ToToken:   r.ToToken,
		Body:      r.Body,
		ServerName:gateway.ServerName,
		Topic: gateway.Topic,
		Address: gateway.ImAddress,
	}); err != nil {
		return nil, PublishMessageErr
	}
	//发送消息结束

	return &SendResponse{}, nil
}

func (l *GatewayLogic) GetImAddress(r *GetImAddressRequest) (*GetImAddressResponse, error) {
	user, err := l.userRpcModel.FindByToken(context.TODO(), &userPb.FindByTokenRequest{Token: r.Token})
	if err != nil {
		return nil, UserNotFoundErr
	}

	len := len(l.imAddressList)
	if len == 0 {
		return nil, ImAddressErr
	}
	index := user.Id % int64(len)

	imConf := l.imAddressList[index]

	if err := l.gatewayModel.Insert(&model.Gateway{
		Token:      r.Token,
		ImAddress:  imConf.Address,
		ServerName: imConf.ServerName,
		Topic:      imConf.Topic,
	}); err != nil {
		return nil, AddGatewayErr
	}

	return &GetImAddressResponse{
		Address: imConf.Address,
	}, nil
}

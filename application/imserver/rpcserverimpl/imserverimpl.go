package rpcserverimpl

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro/broker"
	"micro-me/application/common/baseerror"
	imPb "micro-me/application/imserver/protos"
	"micro-me/application/imserver/service"
)

type(
	ImRpcServerImpl struct{
		publisherServiceMap map[string]*service.RabbitMqService
	}
)

var (
	PublishMessageErr = baseerror.NewBaseError("发送消息失败")
)

func NewImRpcServerImpl(publisherServiceMap map[string]*service.RabbitMqService)*ImRpcServerImpl{
	return &ImRpcServerImpl{publisherServiceMap:publisherServiceMap}
}


func (impl *ImRpcServerImpl)PublishMessage(ctx context.Context, req *imPb.PublishMessageRequest, rsp *imPb.PublishMessageResponse) error{

	body, err := json.Marshal(req)
	if err != nil {
		return PublishMessageErr
	}

	key := req.ServerName + req.Topic
	publisher,ok  := impl.publisherServiceMap[key]
	if !ok {
		return baseerror.NewBaseError("不存在的推送服务")
	}

	publisher.Publisher(&broker.Message{
		Body: body,
	})

	return nil
}

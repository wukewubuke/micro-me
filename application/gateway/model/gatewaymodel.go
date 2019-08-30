package model

import (
	"github.com/go-xorm/xorm"
	"time"
)

type(

	Gateway struct {
		Id int64 `json:"id"`
		Token string `json:"token" xorm:"varchar(100) notnull 'token'"`
		ImAddress string `json:"imAddress" xorm:"varchar(60) notnull 'im_address'"`
		ServerName string `json:"server_name" xorm:"varchar(60) notnull 'server_name'"`
		Topic string `json:"topic" xorm:"varchar(60) notnull 'topic'"`
		CreateTime time.Time `json:"create_time" xorm:"DateTime 'create_time'"`
		UpdateTime time.Time `json:"update_time" xorm:"DateTime 'update_time'"`
	}


	GateWayModel struct{
		mysql *xorm.Engine
	}
)


func NewGatewayModel(mysql *xorm.Engine) *GateWayModel{
	return &GateWayModel{
		mysql:mysql,
	}
}

func (m *GateWayModel)Insert(gateway *Gateway)error{

	has, err := m.FindByServerNameAddressTokenTopic(gateway.ServerName,
		gateway.ImAddress, gateway.Token, gateway.Topic)
	if err == nil && has != nil && has.Id >0  {
		return nil
	}


	if _, err := m.mysql.Insert(gateway); err != nil {
		return err
	}
	return nil
}

func (m *GateWayModel)FindByToken(token string) (*Gateway, error){
	gateway := new(Gateway)
	if _, err := m.mysql.Where("token = ?", token).Get(gateway); err != nil {
		return nil, err
	}
	return gateway, nil
}

func (m *GateWayModel)FindByServerNameAddressTokenTopic(serverName, address, token, topic string) (*Gateway, error){
	gateway := new(Gateway)
	if _, err := m.mysql.Where("token = ? and server_name = ? and im_address = ? and topic = ?",
		token, serverName, address, topic).Get(gateway); err != nil {
		return nil, err
	}
	return gateway, nil
}



func (m *GateWayModel)FindByImAddress(imAddress string) ([]*Gateway, error){
	gateways := []*Gateway(nil)
	if err := m.mysql.Where("im_address = ?", imAddress).Find(&gateways); err != nil {
		return nil, err
	}
	return gateways, nil
}

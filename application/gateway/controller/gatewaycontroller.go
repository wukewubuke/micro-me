package controller

import (
	"github.com/gin-gonic/gin"
	"micro-me/application/common/baseresponse"
	"micro-me/application/gateway/logic"
)

type(
	GatewayController struct {
		gatewayLogic *logic.GatewayLogic
	}



)


func NewGatewayController(gatewayLogic *logic.GatewayLogic)*GatewayController{
	return &GatewayController{
		gatewayLogic:gatewayLogic,
	}
}


func (c *GatewayController)Send(context *gin.Context){
	r := new(logic.SendRequest)
	if err := context.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(context, err)
		return
	}
	res, err := c.gatewayLogic.Send(r)
	baseresponse.HttpResponse(context, res, err)
}


func (c *GatewayController)GetImAddress(context *gin.Context){
	r := new(logic.GetImAddressRequest)
	if err := context.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(context, err)
		return
	}
	res, err := c.gatewayLogic.GetImAddress(r)
	baseresponse.HttpResponse(context, res, err)
}


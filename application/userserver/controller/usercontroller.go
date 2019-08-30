package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"micro-me/application/common/baseresponse"
	"micro-me/application/userserver/logic"
)

type(
	UserController struct {
		userLogic *logic.UserLogic
	}
)


func NewUserController(userLogic *logic.UserLogic) *UserController{
	return &UserController{
		userLogic: userLogic,
	}
}


func (c *UserController)Login(context *gin.Context) {
	r := new(logic.LoginRequest)

	fmt.Printf("%+v\n",r)
	if err := context.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(context, err)
		return
	}
	res, err := c.userLogic.Login(r)
	baseresponse.HttpResponse(context, res, err)
}


func (c *UserController)Register(context *gin.Context) {
	r := new(logic.RegisterRequest)

	fmt.Printf("%+v\n",r)
	if err := context.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(context, err)
		return
	}
	res, err := c.userLogic.Register(r)
	baseresponse.HttpResponse(context, res, err)
}

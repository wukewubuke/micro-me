package logic

import (
	"crypto/md5"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"micro-me/application/common/baseerror"
	"micro-me/application/common/middleware"
	"micro-me/application/userserver/model"
	"time"
)

type(
	UserLogic struct {
		userModel *model.MembersModel
	}

	LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	LoginResponse struct {
		Token string `json:"token"`
		AccessToken string `json:"accessToken"`
		ExpireAt int64 `json:"expireAt"`
		TimeStamp int64 `json:"timestamp"`
	}


	RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	RegisterResponse struct {
	}
)



var (
	NotFoundUser = baseerror.NewBaseError("用户名不存在")
	UserOrPassErr = baseerror.NewBaseError("用户名或密码错误")
	AccessTokenErr = baseerror.NewBaseError("生成签名错误")
	CreateMemberErr = baseerror.NewBaseError("注册失败")
)


func NewUserLogic(userModel *model.MembersModel)*UserLogic {
	return &UserLogic{
		userModel:userModel,
	}
}


func (l *UserLogic)Login(r *LoginRequest)(*LoginResponse, error) {
	user, err := l.userModel.FindByUsername(r.Username)
	if err != nil {
		return nil, NotFoundUser
	}
	if user.Password != fmt.Sprintf("%x",md5.Sum([]byte(r.Password))) {
		return  nil, UserOrPassErr
	}
	expired := time.Now().Add(148 *time.Hour).Unix()
	accessToken, err := l.createAccessToken(expired)
	if err != nil {
		return  nil, AccessTokenErr
	}
	return &LoginResponse{
		Token:user.Token,
		AccessToken: accessToken,
		ExpireAt:expired,
		TimeStamp:time.Now().Unix(),

	}, nil
}


func (l *UserLogic)Register(r *RegisterRequest)(*RegisterResponse,error){
	member := &model.Members{
		Token: uuid.NewV4().String(),
		Username: r.Username,
		Password:fmt.Sprintf("%x",md5.Sum([]byte(r.Password))),
	}
	if err := l.userModel.InsertMember(member); err != nil {
		return nil, CreateMemberErr
	}

	return &RegisterResponse{}, nil
}


func(l *UserLogic)createAccessToken(expired int64)(string, error){
	claims := &jwt.StandardClaims{
		ExpiresAt: expired,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	accessToken, err := token.SignedString([]byte(middleware.UserSignKey))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}
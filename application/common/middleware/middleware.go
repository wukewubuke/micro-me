package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"micro-me/application/common/baseerror"
	"micro-me/application/common/baseresponse"
)

const (
	DefaultField = "Authorization"
	UserSignKey = "admin.sign"
)


var (
	AccessTokenVaildErr = baseerror.NewBaseError("鉴权失败")
	AccessTokenValidationErrorExpired = baseerror.NewBaseError("Token过期")
	AccessTokenValidationErrorMalformed = baseerror.NewBaseError("Token格式错误")
)
func ValidAccessToken(context *gin.Context)  {
	authorization := context.GetHeader(DefaultField)
	token,err := jwt.Parse(authorization, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(UserSignKey),nil
	})
	if err != nil {
		if err ,ok := err.(*jwt.ValidationError);ok {

			if err.Errors & jwt.ValidationErrorMalformed != 0 {
				baseresponse.HttpResponse(context,nil, AccessTokenValidationErrorMalformed)
				context.Abort()
				return
			}
			if err.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
				baseresponse.HttpResponse(context,nil, AccessTokenValidationErrorExpired)
				context.Abort()
				return
			}
		}
		baseresponse.HttpResponse(context,nil, AccessTokenVaildErr)
		context.Abort()
		return
	}
	if token.Valid {
		context.Next()
	}else{
		baseresponse.HttpResponse(context,nil, AccessTokenVaildErr)
		context.Abort()
	}

}

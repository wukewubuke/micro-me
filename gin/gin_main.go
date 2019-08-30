package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	proto "micro-me/proto"
)





type Greeter struct{
	Tag string
}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name + " " + g.Tag
	return nil
}


type (
	Config struct {
		Version string
		Greeter struct {
			Name string
		}
		Etcd struct {
			Addrs []string
			UserName string
			Password string
		}
	}

	User struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Token    struct {
			AccessToken string `json:"accessToken"`
			ExpiresAt   int64  `json:"expiresAt"`
			Timestamp   int64  `json:"timestamp"`
		}
	}


)




func main() {


	configFile := flag.String("f","./config/config_rpc.json","please use config_rpc.json")
	conf := new(Config)

	if err := config.LoadFile(*configFile); err != nil {
		log.Fatal(err)
	}


	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}




	//log.Infof("%+v",conf)
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = conf.Etcd.Addrs
		})




	// Create a new service. Optionally include some options here.
	service := web.NewService(
		web.Name("gin.api.service"),
		web.Registry(etcdRegisty),
		web.Address(":8080"),
	)


	router := gin.Default()
	user := &User{
		Name: "admin",
		Password: "123456",
	}
	router.GET("/user/login", func(context *gin.Context) {
		username := context.Query("username")
		passwrod := context.Query("passwrod")
		fmt.Println(username, passwrod)
		if passwrod != user.Password || username != user.Name {
			context.JSON(200, "pwd err")
			return
		}
		claims := &jwt.StandardClaims{
			ExpiresAt:time.Now().Add(30*time.Second).Unix(),
		}
		expired := time.Now().Add(148 * time.Hour).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
		accessToken, err := token.SignedString([]byte("vector.sign"))
		if err != nil {
			context.JSON(200, "accessToken err")
			return
		}
		user.Token.ExpiresAt = expired
		user.Token.AccessToken = accessToken
		user.Token.Timestamp = time.Now().Unix()

		context.JSON(200, user)
	})
	authorizationRouter := router.Group("/user")
	authorizationRouter.Use(ValidAccessToken)
	authorizationRouter.POST("/user/list", func(context *gin.Context) {

		context.JSON(200,"ok")
		return
	})
	service.Handle("/", router)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}


func ValidAccessToken(context *gin.Context)  {
	authorization := context.GetHeader("Authorization")
	log.Println(authorization)
	token,err := jwt.Parse(authorization, func(token *jwt.Token) (i interface{}, e error) {
		return []byte("vector.sign"),nil
	})
	if err != nil {
		if err ,ok := err.(*jwt.ValidationError);ok {

			if err.Errors & jwt.ValidationErrorMalformed != 0 {
				context.JSON(200,"ValidationErrorMalformed")
				return
			}
			if err.Errors & (jwt.ValidationErrorExpired | jwt.ValidationErrorNotValidYet) != 0 {
				context.JSON(200,"ValidationErrorExpired")
				return
			}
		}
		context.JSON(200,"ValidationError")
		return
	}
	if token.Valid {
		context.Next()
	}else{
		context.JSON(200,"no Valid ")
	}

}




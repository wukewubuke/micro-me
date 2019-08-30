package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/micro/cli"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	userApiConfig "micro-me/application/userserver/cmd/config"
	"micro-me/application/userserver/controller"
	"micro-me/application/userserver/logic"
	"micro-me/application/userserver/model"
)

func main(){
	userApiFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config.json",
		Usage: "please use xxx -f config.json",
	}

	configFile := flag.String(userApiFlag.Name, userApiFlag.Value, userApiFlag.Usage)
	flag.Parse()
	conf := new(userApiConfig.UserApiConfig)

	if err := config.LoadFile(*configFile); err != nil {
		log.Fatal(err)
	}

	if err := config.Scan(conf); err != nil {
		log.Fatal(err)
	}

	//log.Printf("%+v",conf)

	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = conf.Etcd.Addrs
		})

	service := web.NewService(
		web.Name(conf.Server.Name),
		web.Registry(etcdRegisty),
		web.Version(conf.Version),
		web.Flags(userApiFlag),
		web.Address(conf.Port),
	)

	mysql, err := xorm.NewEngine(conf.Engine.Name, conf.Engine.DataSource)
	if err != nil {
		log.Fatal(err)
	}
	userModel := model.NewMembersModel(mysql)
	userLogic := logic.NewUserLogic(userModel)
	userController := controller.NewUserController(userLogic)

	// Init will parse the command line flags.
	router := gin.Default()

	
	userRouterGroup := router.Group("/user")
	{
		userRouterGroup.POST("/login", userController.Login)
		userRouterGroup.POST("/register", userController.Register)
	}



	service.Handle("/",router)

	// Run the service
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

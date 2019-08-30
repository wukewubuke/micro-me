package main

import (
	"flag"
	"micro-me/application/common/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-micro/web"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"log"
	gatewayConfig "micro-me/application/gateway/cmd/config"
	"micro-me/application/gateway/controller"
	"micro-me/application/gateway/logic"
	"micro-me/application/gateway/model"
	imPb "micro-me/application/imserver/protos"
	userPb "micro-me/application/userserver/protos"
)

func main() {
	gatewayApiFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config.json",
		Usage: "please use xxx -f config.json",
	}

	configFile := flag.String(gatewayApiFlag.Name, gatewayApiFlag.Value, gatewayApiFlag.Usage)
	flag.Parse()
	conf := new(gatewayConfig.GatewayApiConfig)

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

	rpcService := micro.NewService(
		micro.Name(conf.UserRpcServer.ClientName),
		micro.Registry(etcdRegisty),
		micro.Flags(gatewayApiFlag),
		micro.Transport(grpc.NewTransport()),
		//客户端熔断降级
		micro.WrapClient(hystrix.NewClientWrapper()),

	)

	rpcService.Init()
	userRpcModel := userPb.NewUserService(conf.UserRpcServer.ServerName, rpcService.Client())


	imRpcModel := imPb.NewImService(conf.ImRpcServer.ServerName, rpcService.Client())

	mysql, err := xorm.NewEngine(conf.Engine.Name, conf.Engine.DataSource)
	if err != nil {
		log.Fatal(err)
	}

	gatewayModel := model.NewGatewayModel(mysql)

	gatewayLogic := logic.NewGatewayLogic(userRpcModel, gatewayModel, conf.ImRpcServer.ImServerList, imRpcModel)
	gatewayController := controller.NewGatewayController(gatewayLogic)

	service := web.NewService(
		web.Name(conf.Server.Name),
		web.Registry(etcdRegisty),
		web.Version(conf.Version),
		web.Flags(gatewayApiFlag),
		web.Address(conf.Port),
	)

	// Init will parse the command line flags.
	router := gin.Default()

	gatewayRouterGroup := router.Group("/gateway")
	gatewayRouterGroup.Use(middleware.ValidAccessToken)
	{
		gatewayRouterGroup.POST("/send", gatewayController.Send)
		gatewayRouterGroup.POST("/address", gatewayController.GetImAddress)

	}

	service.Handle("/", router)

	// Run the service
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

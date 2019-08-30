package main

import (
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	rl "github.com/juju/ratelimit"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
	"log"
	userRpcconfig "micro-me/application/userserver/cmd/config"
	"micro-me/application/userserver/model"
	"micro-me/application/userserver/rpcserverimpl"

	userPb "micro-me/application/userserver/protos"
)

func main() {

	userRpcFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config_rpc.json",
		Usage: "please use xxx -f config_rpc.json",
	}

	configFile := flag.String(userRpcFlag.Name, userRpcFlag.Value, userRpcFlag.Usage)
	flag.Parse()
	conf := new(userRpcconfig.UserRpcConfig)

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

	//限流
	b := rl.NewBucketWithRate(float64(conf.Server.RateLimit), int64(conf.Server.RateLimit))
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Version(conf.Version),
		micro.Transport(grpc.NewTransport()),
		//限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(b, false)),
		micro.Flags(userRpcFlag),
	)

	// Init will parse the command line flags.
	service.Init()

	mysql, err := xorm.NewEngine(conf.Engine.Name, conf.Engine.DataSource)
	if err != nil {
		log.Fatal(err)
	}

	membersModel := model.NewMembersModel(mysql)
	userRpcServer := rpcserverimpl.NewUserRpcServer(membersModel)

	if err := userPb.RegisterUserHandler(service.Server(), userRpcServer); err != nil {
		log.Fatal(err)
	}

	// Run the service
	if err := service.Run(); err != nil {
		log.Println(err)
	}

}

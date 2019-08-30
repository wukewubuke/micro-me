package main

import (
	"flag"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/transport/grpc"
	rabbitmq2 "github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	imConfig "micro-me/application/imserver/cmd/config"
	imPb "micro-me/application/imserver/protos"
	"micro-me/application/imserver/rpcserverimpl"
	imservice "micro-me/application/imserver/service"
)

func main() {
	imFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config_rpc.json",
		Usage: "please use xxx -f config_rpc.json",
	}

	configFile := flag.String(imFlag.Name, imFlag.Value, imFlag.Usage)
	flag.Parse()
	conf := new(imConfig.ImRpcConfig)

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
	log.Println(conf)

	service := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Version(conf.Version),
		micro.Flags(imFlag),
		micro.Transport(grpc.NewTransport()),
	)

	service.Init()

	publisherServerMap := make(map[string]*imservice.RabbitMqService)
	for _, item := range conf.ImServerList {
		amqpAddress := item.AmqpAddress
		p, err := imservice.NewRabbitMqService(item.Topic,
			rabbitmq2.NewBroker(func(options *broker.Options) {
				options.Addrs = amqpAddress
			}))
		if err != nil {
			log.Fatal(err)
		}
		publisherServerMap[item.ServerName+item.Topic] = p
	}

	imRpcServer := rpcserverimpl.NewImRpcServerImpl(publisherServerMap)
	if err := imPb.RegisterImHandler(service.Server(), imRpcServer); err != nil {
		log.Fatal(err)
	}

	// Run the service
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

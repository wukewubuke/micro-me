package main

import (
	"flag"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	rabbitmq2 "github.com/micro/go-plugins/broker/rabbitmq"
	"github.com/micro/go-plugins/registry/etcdv3"
	"log"
	imConfig "micro-me/application/imserver/cmd/config"
	imservice "micro-me/application/imserver/service"
)

func main(){
	imFlag := cli.StringFlag{
		Name:  "f",
		Value: "./config/config.json",
		Usage: "please use xxx -f config.json",
	}

	configFile := flag.String(imFlag.Name, imFlag.Value, imFlag.Usage)
	flag.Parse()
	conf := new(imConfig.ImConfig)

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
	rabbitMqBroker := rabbitmq2.NewBroker(func(options *broker.Options) {
		options.Addrs = conf.RabbitMq.Addresses
	})

	service := micro.NewService(
		micro.Name(conf.Server.Name),
		micro.Registry(etcdRegisty),
		micro.Version(conf.Version),
		micro.Flags(imFlag),
	)

	service.Init()
	rabbitMqService, err := imservice.NewRabbitMqService(conf.RabbitMq.Topic, rabbitMqBroker)
	if err != nil {
		log.Fatal(err)
	}

	imService, err := imservice.NewImService(conf.RabbitMq.Topic, rabbitMqService,func(service *imservice.ImService) {
		service.Address = conf.Port
	})
	if err != nil {
		log.Fatal(err)
	}




	go imService.Subscribe()
	go imService.Run()
	// Run the service
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}

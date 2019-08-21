package main

import (
	"context"
	"flag"
	"fmt"
	rl "github.com/juju/ratelimit"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	"github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
	"github.com/prometheus/common/log"
	"math/rand"
	"strconv"
	"time"

	"github.com/micro/go-micro"
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
)




func main() {


	configFile := flag.String("f","./config/config.json","please use config.json")
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

	//限流
	limit := 2
	b := rl.NewBucketWithRate(float64(limit), int64(limit))



	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name(conf.Greeter.Name),
		micro.Registry(etcdRegisty),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(b, false)),
	)

	// Init will parse the command line flags.
	service.Init()

	greeter := &Greeter{
		Tag: strconv.Itoa(rand.Int()),
	}

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), greeter)



	// 初始化
	if err := broker.Init(); err != nil {
		log.Fatal(err)
	}
	if err := broker.Connect(); err != nil {
		log.Fatal(err)
	}
	go publisher()
	go subscribe()



	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

//消息发布与订阅
var topic = "demo.topic"

func publisher() {
	t := time.NewTicker(time.Second)
	for e := range t.C {
		msg := &broker.Message{
			Header: map[string]string{
				"Tag": strconv.Itoa(rand.Int()),
			},
			Body: []byte(e.String()),
		}
		if err := broker.Publish(topic, msg); err != nil {

			log.Info("[publisher err] : %+v", err)
		}
	}
}

func subscribe() {

	if _, err := broker.Subscribe(topic, func(event broker.Event) error {
		fmt.Printf("subscribe received msg : %s,Header is %+v",
			string(event.Message().Body),
			event.Message().Header,
		)
		fmt.Println()
		return nil
	}); err != nil {
		log.Info("[subscribe err] : %+v", err)
	}
}

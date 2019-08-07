package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"

	"github.com/micro/go-micro"
	proto "micro-me/proto"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = []string{"127.0.0.1:2379"}
		})

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Registry(etcdRegisty),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

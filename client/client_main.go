package main

import (
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
	proto "micro-me/proto"
	"time"
)


//客户端熔断与降级
type MyClientWrapper struct {
	client.Client
}

func (c *MyClientWrapper)Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(e error) error {
		//备用服务就是服务降级，如果主服务错误，则调用备用服务
		fmt.Println("这是一个备用的服务")
		return nil
	})
}


// NewClientWrapper returns a hystrix client Wrapper.
func NewMyClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &MyClientWrapper{c}
	}
}

func main() {

	etcdRegisty := etcdv3.NewRegistry(
		func(options *registry.Options) {
			options.Addrs = []string{"127.0.0.1:2379"}
		})

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.Registry(etcdRegisty),
		//客户端熔断降级
		micro.WrapClient(NewMyClientWrapper()),

	)
	service.Init()

	// Create new greeter client
	greeter := proto.NewGreeterService("greeter", service.Client())





	t:= time.NewTicker(100 * time.Millisecond)
	for e:=range t.C {

		// Call the greeter
		rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "John"})

		//服务器限流每秒钟2次请求，如果超过请求则返回错误
		if err != nil {
			fmt.Printf("err ===> %v, [%v]\n",err,e)
		} else {
			fmt.Printf("msg ===> %v, [%v]\n",rsp.Greeting, e)
		}
	}


}

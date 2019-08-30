package config

import "micro-me/application/common/config"

type (
	GatewayApiConfig struct {
		Version string
		Port    string
		Server  struct {
			Name      string
			RateLimit int
		}
		Etcd struct {
			Addrs    []string
			UserName string
			Password string
		}
		Engine struct {
			Name       string
			DataSource string
		}
		UserRpcServer *config.UserRpcServer

		ImRpcServer struct{
			ServerName string
			ImServerList []*config.ImRpcServer
		}
	}
)

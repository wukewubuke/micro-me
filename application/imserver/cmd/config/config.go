package config

import "micro-me/application/common/config"

type (
	ImConfig struct {
		Version string
		Port string
		Server  struct {
			Name      string
			RateLimit int
		}
		Etcd struct {
			Addrs    []string
			UserName string
			Password string
		}

		RabbitMq *config.RabbitMq
	}


	ImRpcConfig struct {
		Version string
		Server  struct {
			Name      string
			RateLimit int
		}
		Etcd struct {
			Addrs    []string
			UserName string
			Password string
		}

		ImServerList []*config.ImRpcServer
	}
)

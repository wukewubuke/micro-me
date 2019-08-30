package config

type (
	UserRpcServer struct {
		ClientName string
		ServerName string
	}
	RabbitMq struct {
		Addresses []string
		Topic     string
	}
	ImRpcServer struct {
		Address     string //im server的ip和端口
		AmqpAddress []string
		Topic       string
		ServerName  string
	}
)

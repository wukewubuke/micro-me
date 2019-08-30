package config

type (
	UserRpcConfig struct {
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
		Engine struct {
			Name       string
			DataSource string
		}
	}
	UserApiConfig struct {
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
	}
)

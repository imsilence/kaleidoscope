package config

type Backend struct {
	Type string
	Host string
	Port int
}

type Authentication struct {
	Type string
}

type ServiceConfig struct {
	Port           int
	Type           string
	Backend        Backend
	Authentication Authentication
}

type ProxyConfig struct {
	Addr     string
	Services []ServiceConfig
}

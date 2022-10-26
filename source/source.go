package source

import (
	"fmt"
)

type (
	Source interface {
		Config() Config
		Exec() []Proxy
	}
	Proxy struct {
		Ip   string
		Port int
	}
	Config struct {
		Name               string
		ExecRate           int
		CheckIpPortRate    int
		CheckIpPortTimeout int
	}
)

func (p Proxy) IpPort() string {
	return fmt.Sprintf("%s:%d", p.Ip, p.Port)
}

func (p Proxy) IsNull() bool {
	return p.Ip == "" || p.Port == 0
}

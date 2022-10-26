package ip_proxy_pool

import (
	"fmt"
	"testing"
	"time"

	"github.com/huizhang-open-source/ip-proxy-pool/source"
)

type CustomSource1 struct {
}

func (CustomSource1) Config() source.Config {
	return source.Config{
		Name:               "CustomSource1",
		ExecRate:           60,
		CheckIpPortTimeout: 3,
	}
}

func (CustomSource1) Exec() []source.Proxy {
	return nil
}

func TestStart(t *testing.T) {
	GetServer().
		RegisterSources([]source.Source{CustomSource1{}}).
		Start()

	for true {
		total := GetServer().GetProxyPoolTotal()
		fmt.Println(total)

		proxy := GetServer().RandomGetOneProxy()
		fmt.Println(proxy)
		time.Sleep(time.Second * 2)
	}
}

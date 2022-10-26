# ip代理池

## 简介

单机版ip代理池

## 特点

- 自动从第三方源抓取ip
- 支持用户注册第三方源
- 支持自动剔除不可用ip

## 使用样例

````go
package main

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
	// 从第三方源抓ip
	return nil
}

func main() {
	
	// 服务启动
	GetServer().
		RegisterSources([]source.Source{
			CustomSource1{},
		}).
		Start()

	for true {
		// 获取池中代理总数
		total := GetServer().GetProxyPoolTotal()
		fmt.Println(total)

		// 随机从池中获取代理
		proxy := GetServer().RandomGetOneProxy()
		fmt.Println(proxy)
		time.Sleep(time.Second * 2)
	}
}
````
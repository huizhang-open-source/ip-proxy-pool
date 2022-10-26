package ip_proxy_pool

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/huizhang-open-source/ip-proxy-pool/source"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type Server struct {
	sources   map[string]source.Source
	proxyPool sync.Map
}

var server *Server

func init() {
	server = &Server{
		sources:   map[string]source.Source{},
		proxyPool: sync.Map{},
	}

	sourceIp66 := source.Ip66{}
	server.sources[sourceIp66.Config().Name] = sourceIp66

	sourceIp89 := source.Ip89{}
	server.sources[sourceIp89.Config().Name] = sourceIp89
}

func GetServer() *Server {
	return server
}

func (s *Server) RegisterSources(sources []source.Source) *Server {
	for _, item := range sources {
		s.sources[item.Config().Name] = item
	}
	return s
}

func (s *Server) Start() {
	for _, sourceItem := range s.sources {
		s.registeredProxy(sourceItem)

		s.monitorProxyHealth(sourceItem)
	}
}

func (s *Server) registeredProxy(sourceItem source.Source) {
	go func(sourceItem source.Source) {
		for true {
			proxyItems := sourceItem.Exec()

			s.registeredProxyToPool(sourceItem, proxyItems)

			time.Sleep(time.Duration(sourceItem.Config().ExecRate) * time.Second)
		}
	}(sourceItem)
}

func (s *Server) monitorProxyHealth(sourceItem source.Source) {
	go func(sourceItem source.Source) {
		for true {
			s.proxyPool.Range(func(key, value interface{}) bool {
				proxy := value.(source.Proxy)
				proxyKey := fmt.Sprintf("%s_%s_%d", sourceItem.Config().Name, proxy.Ip, proxy.Port)
				if cast.ToString(key) == proxyKey {
					if !s.ipPortIsEffect(proxy.Ip, proxy.Port, sourceItem.Config().CheckIpPortTimeout) {
						s.proxyPool.Delete(proxyKey)

						logrus.WithFields(logrus.Fields{
							"source_name": sourceItem.Config().Name,
							"ip":          proxy.Ip,
							"port":        proxy.Port,
						}).Infof("[Server.monitorProxyHealth] Ip:port invalid")
					}
				}
				return true
			})
			time.Sleep(time.Duration(sourceItem.Config().CheckIpPortRate) * time.Second)
		}
	}(sourceItem)
}

func (s *Server) registeredProxyToPool(sourceItem source.Source, proxyItems []source.Proxy) {
	for _, proxyItem := range proxyItems {
		proxyKey := fmt.Sprintf("%s_%s_%d", sourceItem.Config().Name, proxyItem.Ip, proxyItem.Port)

		if _, ok := s.proxyPool.Load(proxyKey); ok {
			continue
		}

		if !s.ipPortIsEffect(proxyItem.Ip, proxyItem.Port, sourceItem.Config().CheckIpPortTimeout) {

			s.proxyPool.Delete(proxyKey)

			logrus.WithFields(logrus.Fields{
				"source_name": sourceItem.Config().Name,
				"ip":          proxyItem.Ip,
				"port":        proxyItem.Port,
			}).Infof("[Server.registeredProxyToPool] Ip:port invalid")
			continue
		}

		if _, ok := s.proxyPool.LoadOrStore(proxyKey, proxyItem); !ok {
			logrus.WithFields(logrus.Fields{
				"source_name": sourceItem.Config().Name,
				"ip":          proxyItem.Ip,
				"port":        proxyItem.Port,
			}).Infof("[Server.registeredProxyToPool] Registered proxy pool success")
		}
	}
}

func (s *Server) GetProxyPoolTotal() int {
	var result int
	s.proxyPool.Range(func(key, value interface{}) bool {
		result++
		return true
	})
	return result
}

func (s *Server) RandomGetOneProxy() source.Proxy {
	var result source.Proxy

	s.proxyPool.Range(func(key, value interface{}) bool {
		result = value.(source.Proxy)
		return false
	})

	return result
}

func (*Server) ipPortIsEffect(ip string, port int, timeout int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Duration(timeout)*time.Second)
	if err != nil {
		return false
	}

	if conn == nil {
		return false
	}

	defer conn.Close()

	return true
}

package source

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

const (
	ip66Name = "ip_66"
	ip66Url  = "http://www.66ip.cn/mo.php?tqsl=100"
)

type Ip66 struct {
}

func (Ip66) Config() Config {
	return Config{
		Name:               ip66Name,
		ExecRate:           60,
		CheckIpPortTimeout: 3,
	}
}

func (Ip66) Exec() []Proxy {
	var exprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	resp, err := http.Get(ip66Url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":         ip66Url,
			"source_name": ip66Name,
		}).WithError(err).Errorf("[Ip66.Exec] http.Get fail")
		return nil
	}

	if resp.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{
			"url":         ip66Url,
			"source_name": ip66Name,
			"StatusCode":  resp.StatusCode,
		}).Errorf("[Ip66.Exec] resp.StatusCode != 200")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	ips := exprIP.FindAllString(string(body), 100)

	var proxyArr []Proxy
	for _, ip := range ips {
		ipStr := strings.TrimSpace(ip)
		ipArr := strings.Split(ipStr, ":")
		if len(ipArr) != 2 {
			continue
		}

		proxyArr = append(proxyArr, Proxy{
			Ip:   ipArr[0],
			Port: cast.ToInt(ipArr[1]),
		})
	}

	return proxyArr
}

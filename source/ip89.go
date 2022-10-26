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
	ip89Name = "ip_89"
	ip89Url  = "https://www.89ip.cn/tqdl.html?api=1&num=200&port=&address=俄罗斯&isp="
)

type Ip89 struct {
}

func (Ip89) Config() Config {
	return Config{
		Name:               ip89Name,
		ExecRate:           60,
		CheckIpPortTimeout: 3,
	}
}

func (Ip89) Exec() []Proxy {
	var exprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	resp, err := http.Get(ip89Url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url":         ip89Url,
			"source_name": ip89Name,
		}).WithError(err).Errorf("[Ip66.Exec] http.Get fail")
		return nil
	}

	if resp.StatusCode != 200 {
		logrus.WithFields(logrus.Fields{
			"url":         ip89Url,
			"source_name": ip89Name,
			"StatusCode":  resp.StatusCode,
		}).Errorf("[Ip66.Exec] resp.StatusCode != 200")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	ips := exprIP.FindAllString(string(body), 200)

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

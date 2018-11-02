package main

import (
	"flag"
	"fmt"
	"httpmq-proxy/common"
	"httpmq-proxy/httpproxy"
	"httpmq-proxy/mqproxy"
)

var (
	proxy = flag.String("proxy", "", "node role type. example(master or slave)")

	mq      = flag.String("mq", "gaozh:gaozh@192.168.1.242:5672", "mq address. example(guest:guest@192.168.1.242:5672)")
	cluster = flag.String("cluster", "", "local cluster api address is must when proxy is slave. example(192.168.1.1:6443)")
	listen  = flag.String("listen", "0.0.0.0:6000", "local http listen address. default(0.0.0.0:6000)")
)

func main() {
	flag.Parse()
	if *proxy == "master" {
		httpHandler, err := httpProxy.NewHttpHandler(*listen, *cluster, *proxy, *mq)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = httpHandler.HttpListen(); err != nil {
			fmt.Println(err)
			return
		}
	}

	if *proxy == "slave" && *cluster != "" {
		mqHandler, err := mqproxy.NewMqHandle(*mq, *proxy)
		if err != nil {
			fmt.Println(err)
			return
		}
		httpSender, err := common.NewHttpSend(*cluster)
		if err != nil {
			fmt.Println(err)
			return
		}
		mqHandler.RecvDataFromQueue(httpSender)
	}

	flag.PrintDefaults()
}

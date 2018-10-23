package main

import (
	"flag"
	"fmt"
	"httpmq-proxy/common"
	"httpmq-proxy/httpproxy"
	"httpmq-proxy/mqproxy"
)

var (
	role = flag.String("role", "", "node role type. example(tce or tde)")

	mq      = flag.String("mq", "gaozh:gaozh@192.168.1.242:5672", "mq address. example(guest:guest@192.168.1.242:5672)")
	apihost = flag.String("apihost", "", "local cluster api address is must when role is tde. example(192.168.1.1:6443)")
	listen  = flag.String("listen", "0.0.0.0:6444", "local http listen address. default(0.0.0.0:6444)")
)

func main() {
	flag.Parse()
	if *role == "tce" {
		httpHandler, err := httpProxy.NewHttpHandler(*listen, *apihost, *role, *mq)
		if err != nil {
			fmt.Println(err)
			return
		}
		httpHandler.HttpListen()
	}

	if *role == "tde" && *apihost != "" {
		mqHandler, err := mqproxy.NewMqHandle(*mq, *role)
		if err != nil {
			fmt.Println(err)
			return
		}
		httpSender, err := common.NewHttpSend(*apihost)
		if err != nil {
			fmt.Println(err)
			return
		}
		mqHandler.RecvDataFromQueue(httpSender)
	}

	flag.PrintDefaults()
}

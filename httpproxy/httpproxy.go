package httpProxy

import (
	"fmt"
	"httpmq-proxy/common"
	"httpmq-proxy/mqproxy"
	"io/ioutil"
	"net/http"
	"time"
)

type httpHandle struct {
	Apihost       string
	ListenAddress string
	MqHandler     mqproxy.MqHandler
}

type httpHandler interface {
	HttpListen()
}

func NewHttpHandler(listenAddress, apihost, role, mqAddress string) (httpHandler, error) {
	mqHandler, err := mqproxy.NewMqHandle(mqAddress, role)
	if err != nil {
		return nil, err
	}
	return &httpHandle{
		Apihost:       apihost,
		ListenAddress: listenAddress,
		MqHandler:     mqHandler,
	}, nil
}

func (self httpHandle) handle(res http.ResponseWriter, req *http.Request) {
	msg := fmt.Sprintf("[%s] %s", req.Method, req.URL)
	body, _ := ioutil.ReadAll(req.Body)
	requestPackage := common.RequestPackage{
		URL:    req.URL,
		Method: req.Method,
		Header: req.Header,
		Body:   body,
	}

	//1、use https local request k8s
	httpSender, err := common.NewHttpSend(self.Apihost)
	if err != nil {
		return
	}
	responPackage, err := httpSender.SendHttpRequest(requestPackage)
	if err != nil {
		return
	}

	/*
		//2、use mq remote request k8s
		requestPackageByte, _ := json.Marshal(requestPackage)
		responByte, err := self.MqHandler.SendDataToQueue(requestPackageByte)
		if err != nil {
			fmt.Printf("%s %v\n", msg, err)
			return
		}

		var responPackage common.ResponsePackage
		if err = json.Unmarshal(responByte, &responPackage); err != nil {
			fmt.Println(err)
			return
		}
	*/

	for key, value := range responPackage.HeadMap {
		res.Header().Set(key, value)
	}

	fmt.Printf("%s %v\n", msg, responPackage.StatusCode)
	res.Write(responPackage.Body)
}

func (self httpHandle) HttpListen() {
	fmt.Printf("listen http address: %s\n", self.ListenAddress)
	server := &http.Server{
		Addr:         self.ListenAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.HandleFunc("/", self.handle)
	server.ListenAndServe()
}

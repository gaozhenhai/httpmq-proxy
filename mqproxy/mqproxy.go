package mqproxy

import (
	"encoding/json"
	"fmt"
	"httpmq-proxy/common"
	"time"

	"github.com/streadway/amqp"
)

type MqHandle struct {
	Address string
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Queue   amqp.Queue
	Msgs    <-chan amqp.Delivery
}

type MqHandler interface {
	RecvDataFromQueue(httpSender common.HttpSender)
	SendDataToQueue(requestPackage common.RequestPackage) (common.ResponsePackage, error)
}

func NewMqHandle(address, role string) (MqHandler, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s", address))
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	ch.Qos(1, 0, false)

	if role == "tce" {
	} else {
	}
	return &MqHandle{
		Channel: ch,
		Address: address,
	}, nil
}

func (self MqHandle) SendDataToQueue(requestPackage common.RequestPackage) (common.ResponsePackage, error) {
	var responPackage common.ResponsePackage
	queue, err := self.Channel.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return responPackage, err
	}
	msgs, err := self.Channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return responPackage, err
	}

	corrId := common.RandomString(32)
	requestPackageByte, _ := json.Marshal(requestPackage)
	if err := self.Channel.Publish("", "dreamx_queue", false, false, amqp.Publishing{
		CorrelationId: corrId,
		ReplyTo:       queue.Name,
		Body:          requestPackageByte,
	}); err != nil {
		return responPackage, err
	}

	for {
		select {
		case d := <-msgs:
			if corrId == d.CorrelationId {
				err := json.Unmarshal(d.Body, &responPackage)
				return responPackage, err
			}
		case <-time.After(15e9):
			return responPackage, fmt.Errorf("request timeout")
		}
	}
	return responPackage, nil
}

func (self MqHandle) RecvDataFromQueue(httpSender common.HttpSender) {
	var requestPackage common.RequestPackage
	queue, err := self.Channel.QueueDeclare("dreamx_queue", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	msgs, err := self.Channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	forever := make(chan bool)
	for n := 0; n < 3; n++ {
		go func() {
			for d := range msgs {
				if err := json.Unmarshal(d.Body, &requestPackage); err != nil {
					fmt.Println(err)
				}
				responsePackage, err := httpSender.SendHttpsRequest(requestPackage)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("----> %s %s %v\n", requestPackage.Method, requestPackage.URL, responsePackage.StatusCode)
				body, _ := json.Marshal(responsePackage)
				if err := self.Channel.Publish("", d.ReplyTo, false, false, amqp.Publishing{CorrelationId: d.CorrelationId, Body: body}); err != nil {
					fmt.Println(err)
				}
			}
		}()
	}
	<-forever
}

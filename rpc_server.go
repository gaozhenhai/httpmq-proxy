package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@192.168.1.242:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("dreamx_queue", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	//err = ch.Qos(1, 0, false)
	//failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("recv====> %v\n", string(d.Body))
			temp := "hello world"
			err = ch.Publish("", d.ReplyTo, false, false, amqp.Publishing{CorrelationId: d.CorrelationId, Body: []byte(temp)})
			failOnError(err, "Failed to publish a message")
			//d.Ack(false)
		}
	}()
	<-forever
}

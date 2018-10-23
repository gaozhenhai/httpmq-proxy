package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@192.168.1.242:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	for {
		corrId := randomString(32)
		err = ch.Publish("", "dreamx_queue", false, false, amqp.Publishing{
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte("gaozh"),
		})
		failOnError(err, "Failed to publish a message")

	BREAK:
		for {
			select {
			case d := <-msgs:
				if corrId == d.CorrelationId {
					fmt.Printf("recv----> %v\n", string(d.Body))
					break BREAK
				}
			case <-time.After(3e9):
				fmt.Println("----------> timeout")
				break BREAK
			}
		}
		time.Sleep(1e8)
	}
	return
}

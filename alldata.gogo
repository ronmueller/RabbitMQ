package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type ticker struct {
	ID                int    `json:"tickerId"`
	EventID           int    `json:"eventId"`
	ResourceOperation string `json:"resourceOperation"`
}

type tickerCreate struct {
	ticker
	ResourceName string `json:"resourceName"`
}
type tickerUpdate struct {
	ticker
	ResourceName string `json:"resourceName"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial("amqp://bunny:39GpHGT49d@rabbitmq.staging.tam-cms.com:5672/")
	// conn, err := amqp.Dial("amqp://bunny:39GpHGT49d@rabbitmq.dev.tam-cms.com:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"feed.alldata", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,         // queue name
		"liveticker.#", // routing key
		"liveticker",   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind liveticker exchange")

	err = ch.QueueBind(
		q.Name,       // queue name
		"metadata.#", // routing key
		"metadata",   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind metadata exchange")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"feed", // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)

			// if d.Body != nil {
			// 	var data ticker
			// 	if len(d.Body) > 0 {
			// 		if err := json.Unmarshal(d.Body, &data); err != nil {
			// 			log.Println(err)
			// 		}
			// 	}

			// 	fmt.Printf("%+v\n", data)
			// 	fmt.Println(data.ID)
			// }

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

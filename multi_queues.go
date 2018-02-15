package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// conn, err := amqp.Dial("amqp://bunny:39GpHGT49d@rabbitmq.staging.tam-cms.com:5672/")
	conn, err := amqp.Dial("amqp://bunny:39GpHGT49d@rabbitmq.dev.tam-cms.com:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q1, err := ch.QueueDeclare(
		"feed.liveticker", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue 1")

	q2, err := ch.QueueDeclare(
		"feed.metadata", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue 2")

	err = ch.QueueBind(
		q1.Name,        // queue name
		"liveticker.#", // routing key
		"liveticker",   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue 1")

	err = ch.QueueBind(
		q2.Name,      // queue name
		"metadata.#", // routing key
		"metadata",   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue 2")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs1, err := ch.Consume(
		q1.Name,      // queue
		"liveticker", // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	msgs2, err := ch.Consume(
		q2.Name,    // queue
		"metadata", // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d1 := range msgs1 {
			log.Printf(" [x] liveticker: %s", d1.Body)
			d1.Ack(false)
		}

		for d2 := range msgs2 {
			log.Printf(" [x] metadata: %s", d2.Body)
			d2.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

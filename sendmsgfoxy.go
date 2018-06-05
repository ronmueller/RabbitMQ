package main

import (
	"encoding/json"
	"flag"
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

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// conn, err := amqp.Dial("amqp://bunny:39GpHGT49d@rabbitmq.staging.tam-cms.com:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	blockings := conn.NotifyBlocked(make(chan amqp.Blocking))
	go func() {
		for b := range blockings {
			if b.Active {
				log.Printf("TCP blocked: %q", b.Reason)
			} else {
				log.Printf("TCP unblocked")
			}
		}
	}()

	// messages
	msgcnt := flag.Int("cnt", 10, "Delivers x messages")
	msgstr := flag.String("msg", "", "Deliver this Message")
	msgtype := flag.String("type", "", "Deliver only this type of Message")
	msgheader := flag.String("header", "", "Deliver this Header")
	flag.Parse()

	// open first channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"feed.input", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for msgs := 0; msgs < *msgcnt; msgs++ {
		header := make(amqp.Table)
		if *msgheader != "" {
			err := json.Unmarshal([]byte(*msgheader), &header)
			failOnError(err, "Failed to read header")
		}

		body := bodyFrom(*msgstr, *msgtype, msgs+1)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         []byte(body),
				Headers:      header,
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s to channel 1", body)
	}

}

func bodyFrom(msg string, msgtype string, id int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var s string
	rnd := rand.Intn(4)
	if (len(msg) < 2) || msg == "" {

		if msgtype == "" {
			switch rnd {
			case 0:
				s = fmt.Sprintf("{\"tickerId\":%d,\"resourceName\":\"tester\",\"resourceOperation\":\"create\"}", id)
			case 1:
				s = fmt.Sprintf("{\"tickerId\":%d,\"resourceName\":\"tester\",\"resourceOperation\":\"update\"}", id)
			case 2:
				s = fmt.Sprintf("{\"tickerId\":%d,\"resourceName\":\"tester\",\"resourceOperation\":\"delete\"}", id)
			case 3:
				s = fmt.Sprintf("{\"tickerId\":%d,\"resourceName\":\"tester\",\"resourceOperation\":\"modify\"}", id)
			}
		} else {
			s = fmt.Sprintf("{\"tickerId\":%d,\"resourceName\":\"tester\",\"resourceOperation\":\"%s\"}", id, msgtype)
		}
	} else {
		s = msg
	}
	return s
}

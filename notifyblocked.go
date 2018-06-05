package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Simply logs when the server throttles the TCP connection for publishers

	// Test this by tuning your server to have a low memory watermark:
	// rabbitmqctl set_vm_memory_high_watermark 0.00000001

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("connection.open: %s", err)
	}
	defer conn.Close()

	blockings := conn.NotifyBlocked(make(chan amqp.Blocking))
	go func() {
		log.Println("loooooop")
		for b := range blockings {
			if b.Active {
				log.Printf("TCP blocked: %q", b.Reason)
			} else {
				log.Printf("TCP unblocked")
			}
		}
	}()

}

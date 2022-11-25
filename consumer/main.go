package main

import (
	"log"

	"github.com/kecci/go-rabbitmq/utility"
)

func main() {
	// Init Rabbit channel and queue
	conn, ch, q := utility.InitRabbitMQ()
	defer conn.Close() // Don't forget to close
	defer ch.Close()   // Don't forget to close

	// We set the consumer for the message
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utility.FailOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-make(chan bool)
}

package main

import (
	"log"

	"github.com/kecci/go-rabbitmq/utility"
	"github.com/streadway/amqp"
)

func main() {
	// Init Rabbit channel and queue
	conn, ch, q := utility.InitRabbitMQ()
	defer conn.Close() // Don't forget to close
	defer ch.Close()   // Don't forget to close

	// We set the payload for the message.
	body := "Golang is awesome - Keep Moving Forward!"
	err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	// If there is an error publishing the message, a log will be displayed in the terminal.
	utility.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Congrats, sending message: %s", body)
}

# go-rabbitmq
Example project of Golang using RabbitMQ

- [go-rabbitmq](#go-rabbitmq)
  - [Publishing a Message With Go](#publishing-a-message-with-go)
  - [Reading Messages With Go](#reading-messages-with-go)
  - [Source](#source)


## Publishing a Message With Go
Now that RabbitMQ is ready to Go, let's connect to it and send a message to the queue. But first, we need the amqp library. To install it, run the following command in your terminal: go get github.com/streadway/amqp.

Then, create a filed called sendMessage.go inside the rabbitmq-go directory. Add the following:

```go
package main

import (
	"log"

	"github.com/streadway/amqp"
)

// Here we set the way error messages are displayed in the terminal.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Here we connect to RabbitMQ or send a message if there are any errors connecting.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// We create a Queue to send the message to.
	q, err := ch.QueueDeclare(
		"golang-queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// We set the payload for the message.
	body := "Golang is awesome - Keep Moving Forward!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	// If there is an error publishing the message, a log will be displayed in the terminal.
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Congrats, sending message: %s", body)
}
```

Once we've created the main.go file, let's test sendMessage with the following terminal command: go run sendMessage.go. You should see something like this:

When you head over to http://localhost:15672, you should see golang-queue in the Queues tab along with the number of messages that have been sent to it.

Awesome! Pat yourself on the back. You've just created a queue and sent a message to it. Now let's learn how to read those messages.

## Reading Messages With Go
Create a file called consumer.go. Inside the file, add the following:

```go
package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"golang-queue", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
```

Now that consume.go is created, let's run it with the following command: go run consume.go. The <-forever line at the end of the file means we'll keep listening to the channel for new messages.

We've made a lot of progress so far. But our sending messages functionality isn't dynamic right now, so let's change that.

Add the following code to the top of the sendMessage.go file, right after `func main() {`.

```go
// Let's catch the message from the terminal.
reader := bufio.NewReader(os.Stdin)
fmt.Println("What message do you want to send?")
mPayload, _ := reader.ReadString('\n')
```

Here, we're creating a new Reader to catch the input from the terminal. This makes sending a message more dynamic. Let's test it with `go run sendMessage.go`.

The beauty of using RabbitMQ is that you can send messages using Go, but read them with Python, JavaScript, PHP, Java, and more! Here's everything in action:

You now know how to spin up a docker image, start RabbitMQ, dynamically send messages to the queue, and read messages from the queue. Bravo! You have everything you need to create amazing applications that efficiently send and receive messages.

Find the files we've created in this tutorial right here.

## Source
- https://x-team.com/blog/set-up-rabbitmq-with-docker-compose/
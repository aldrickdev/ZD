package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func checkError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:admin@broker:5672")
	checkError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	checkError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"zendesk", // name
		"topic",   // type
		true,      //durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	checkError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // Queue Name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	checkError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,    // queue name
		"",        // routing key
		"zendesk", // exchange
		false,     // no-wait
		nil,       // arguments
	)
	checkError(err, "Failed to bind a queue")

	routing_key := "new.userEvent"
	err = ch.QueueBind(
		q.Name,      // queue name
		routing_key, // routing key
		"zendesk",   //exchange
		false,
		nil,
	)
	checkError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // Queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	checkError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for message. To exit press CTRL+C")
	<-forever
}

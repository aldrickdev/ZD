package main

import (
	"fmt"
	"log"
	"os"
	"zd/internal/rabbitmq"
	"zd/internal/utils"
)

func checkError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %q", msg, err)
	}
}

func init() {
	// Load the environment variables
	utils.LoadEnvVars()
}

func main() {
	arguments := os.Args[1:]
	if len(arguments) != 1 {
		log.Fatal("You need to provide just one parameter, the name of the exchange routing key you would like to bind too. examples (marquee, userevent)")
	}

	routingKey := arguments[0]

	// Create and configure the RabbitMQ instance ===
	rabbitMQ := rabbitmq.New()
	rmqConnectionString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		utils.Env.RMQ_USER,
		utils.Env.RMQ_PASS,
		utils.Env.RMQ_DOMAIN,
		utils.Env.RMQ_PORT,
	)
	err := rabbitMQ.Connect(rmqConnectionString)
	checkError(err, "Failed to connect to RabbitMQ")

	err = rabbitMQ.DeclareExchange("zendesk", "topic")
	checkError(err, "Failed to declare an exchange")

	queue, err := rabbitMQ.DeclareQueue()
	checkError(err, "Failed to declare a queue")

	err = rabbitMQ.QueueBind(queue, routingKey)
	checkError(err, "Failed to bind a queue")

	msgs, err := rabbitMQ.ConsumeQueue(queue.Name)
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

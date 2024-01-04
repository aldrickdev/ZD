package main

import (
	"fmt"
	"log"
	"time"
	"zd/envvars"
	"zd/internal/core/domain"
	"zd/internal/driven/rabbitmq"
	"zd/internal/utils"
)

func init() {
	envvars.LoadEnvVars()
}

func main() {
	// Creating and Configuring the RabbitMQ Driven Actor
	rabbitMQ := rabbitmq.New()
	rmqConnectionString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		envvars.Env.RMQ_USER,
		envvars.Env.RMQ_PASS,
		envvars.Env.RMQ_DOMAIN,
		envvars.Env.RMQ_PORT,
	)
	fmt.Printf("Connection String for RabbitMQ: %s\n", rmqConnectionString)
	err := rabbitMQ.Connect(rmqConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	err = rabbitMQ.DeclareExchange("zendesk", "topic")
	if err != nil {
		log.Fatalf("Failed to decalre an exchange: %v", err)
	}

	rabbitMQ.RegisterExchangeRoute("marquee")
	rabbitMQ.RegisterExchangeRoute("userevent")

	// Starting Graceful Shutdown Channel
	utils.GracefulShutdown([]utils.Closable{
		rabbitMQ,
	})

	for {
		ue := domain.MarqueeData{
			UserName:  "Aldrick",
			EventName: "TestEvent",
		}

		rabbitMQ.ExchangeRoutes[0].Publish(ue)

		time.Sleep(5 * time.Second)
	}

	// Loop to prevent main routine from stopping
	// var forever chan struct{}
	// <-forever
}

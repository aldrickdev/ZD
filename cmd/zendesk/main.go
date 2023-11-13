package main

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"zd/envvars"
	httpserver "zd/internal/adapters/driver/httpServer"
	"zd/internal/applications/core/zendesk"
	"zd/internal/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	envvars.LoadEnvVars()
}
func checkError(err error, msg string) {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}

func main() {
	ep := zendesk.NewZendeskMock()
	httpServer := httpserver.NewAdapter(*ep)

	// TODO: Use environment variables for username, password, service name and port
	connect, err := amqp.Dial("amqp://admin:admin@broker:5672")
	checkError(err, "Failed to connect to RabbitMQ")
	defer connect.Close()

	ch, err := connect.Channel()
	checkError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"zendesk", // Exchange Name
		"topic",   // Exchange Type
		true,      // Durable
		false,     // auto-delete
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	checkError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	event := zendesk.UserEvent{
		UserID:    1,
		EventName: "test",
		Points:    10,
	}
	data, err := json.Marshal(event)
	checkError(err, "Failed to marshal JSON")

	err = ch.PublishWithContext(
		ctx,
		"zendesk",       //Exchange Name
		"new.userEvent", // Routing Key
		false,           // mandatory
		false,           // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	checkError(err, "Failed to publish")

	utils.GracefuleShutdown()

	httpServer.Run()
}

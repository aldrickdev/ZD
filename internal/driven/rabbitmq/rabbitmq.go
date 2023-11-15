package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zd/internal/core/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	connect      *amqp.Connection
	channel      *amqp.Channel
	ExchangeName string
}

func New() *RabbitMQ {
	return &RabbitMQ{}
}

func (r *RabbitMQ) Connect(connectionString string) error {
	connect, err := amqp.Dial(connectionString)
	if err != nil {
		return err
	}
	r.connect = connect

	ch, err := connect.Channel()
	if err != nil {
		return err
	}

	r.channel = ch

	return nil
}

func (r *RabbitMQ) DeclareExchange(exchangeName, exchangeType string) error {
	err := r.channel.ExchangeDeclare(
		exchangeName, // Exchange Name
		exchangeType, // Exchange Type
		true,         // Durable
		false,        // auto-delete
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	r.ExchangeName = exchangeName
	return nil
}

func (r *RabbitMQ) Publish(ue domain.UserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("Data being published: %v\n", ue)
	fmt.Printf("Published to Exchange: %s\n", r.ExchangeName)

	data, err := json.Marshal(ue)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(
		ctx,
		r.ExchangeName,
		"new.userevent", // Routing Key
		false,           // mandatory
		false,           // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) GracefulShutdown() {
	fmt.Println("Closing Channel and Connection to RabbitMQ")
	r.channel.Close()
	r.connect.Close()
	fmt.Println("Closed Channel and Connection to RabbitMQ")
}

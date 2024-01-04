package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zd/internal/core/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Route struct {
	channel      *amqp.Channel
	exchangeName string
	key          string
}

func NewRoute(routingKey, exchangeName string, channel *amqp.Channel) Route {
	return Route{
		channel:      channel,
		exchangeName: exchangeName,
		key:          routingKey,
	}
}

func (r Route) Publish(data any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("Data being published: %v\n", data)
	fmt.Printf("Published to Exchange: %s\n", r.exchangeName)

	rawData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(
		ctx,
		r.exchangeName,
		r.key, // Routing Key
		false, // mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        rawData,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

type RabbitMQ struct {
	connect        *amqp.Connection
	channel        *amqp.Channel
	exchangeName   string
	ExchangeRoutes []Route
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

	r.exchangeName = exchangeName
	return nil
}

func (r *RabbitMQ) RegisterExchangeRoute(routingKey string) {
	newRoute := NewRoute(routingKey, r.exchangeName, r.channel)

	r.ExchangeRoutes = append(r.ExchangeRoutes, newRoute)
}

func (r *RabbitMQ) Publish(ue domain.UserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("Data being published: %v\n", ue)
	fmt.Printf("Published to Exchange: %s\n", r.exchangeName)

	data, err := json.Marshal(ue)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(
		ctx,
		r.exchangeName,
		"userevent", // Routing Key
		false,       // mandatory
		false,       // Immediate
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

func (r *RabbitMQ) PublishBatch(ue []*domain.UserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	events := []domain.UserEvent{}
	for _, evt := range ue {
		events = append(events, *evt)
	}

	fmt.Printf("Data being published: %v\n", events)

	data, err := json.Marshal(events)
	if err != nil {
		return err
	}

	err = r.channel.PublishWithContext(
		ctx,
		r.exchangeName,
		"marquee", // Routing Key
		false,     // mandatory
		false,     // Immediate
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

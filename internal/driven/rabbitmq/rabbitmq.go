package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zd/internal/core/domain"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RouteTypeUserEventIDData   = "USER_EVENT_ID_DATA"
	RouteTypeUserEventNameData = "USER_EVENT_NAME_DATA"
)

type Route struct {
	channel      *amqp.Channel
	exchangeName string
	key          string
	routeType    string
}

func NewRoute(routingKey, routeDataType, exchangeName string, channel *amqp.Channel) Route {
	return Route{
		channel:      channel,
		exchangeName: exchangeName,
		key:          routingKey,
		routeType:    routeDataType,
	}
}

func (r Route) Publish(data domain.FullUserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("Data being published: %v\n", data)
	fmt.Printf("Published to Exchange: %s\n", r.exchangeName)

	var rawData []byte
	var err error

	switch r.routeType {
	case RouteTypeUserEventIDData:
		refinedData := domain.UserEventIDData{
			UserID:  data.User.ID,
			EventID: data.Event.ID,
		}

		rawData, err = json.Marshal(refinedData)
		if err != nil {
			return err
		}

	case RouteTypeUserEventNameData:
		refinedData := domain.UserEventNameData{
			UserName:  data.User.Name,
			EventName: data.Event.Name,
		}

		rawData, err = json.Marshal(refinedData)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid route type: %s", r.routeType)
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
	ExchangeRoutes map[string]Route
}

func New() *RabbitMQ {
	return &RabbitMQ{
		ExchangeRoutes: map[string]Route{},
	}
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

func (r *RabbitMQ) RegisterExchangeRoute(routingKey, dataType string) Route {
	newRoute := NewRoute(routingKey, dataType, r.exchangeName, r.channel)
	r.ExchangeRoutes[routingKey] = newRoute

	return newRoute
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

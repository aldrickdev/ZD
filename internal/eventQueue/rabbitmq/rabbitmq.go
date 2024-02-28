package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"zd/internal/core"

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

func (r Route) Publish(data *core.FullUserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if data == nil {
		fmt.Println("No Data to publish")
		return nil
	}

	var rawData []byte
	var err error

	switch r.routeType {
	case RouteTypeUserEventIDData:
		refinedData := core.UserEventIdData{
			UserId:  data.User.Id,
			EventId: data.Event.Id,
		}

		rawData, err = json.Marshal(refinedData)
		if err != nil {
			return err
		}

	case RouteTypeUserEventNameData:
		refinedData := core.UserEventNameData{
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

	fmt.Println("Successfully published data to queue")

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
		false,        // Durable
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

func (r *RabbitMQ) DeclareQueue() (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		"",    // Queue Name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *RabbitMQ) QueueBind(queue amqp.Queue, routingKey string) error {
	return r.channel.QueueBind(
		queue.Name,     // queue name
		routingKey,     // routing key
		r.exchangeName, //exchange
		false,          // no wait
		nil,
	)
}

func (r *RabbitMQ) ConsumeQueue(queueName string) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		queueName, // Queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
}

func (r *RabbitMQ) GracefulShutdown() {
	fmt.Println("Closing Channel and Connection to RabbitMQ")
	r.channel.Close()
	r.connect.Close()
	fmt.Println("Closed Channel and Connection to RabbitMQ")
}

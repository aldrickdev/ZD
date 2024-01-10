package dependencies

import (
	"zd/internal/rabbitmq"
)

type QueueBroker interface {
	Connect(string) error
	DeclareExchange(string, string) error
	RegisterExchangeRoute(string, string) rabbitmq.Route
	GracefulShutdown()
}

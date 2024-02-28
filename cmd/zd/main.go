package main

import (
	"fmt"
	"log"
	"strconv"
	rediscache "zd/internal/cache/redis"
	"zd/internal/core"
	"zd/internal/eventQueue/rabbitmq"
	"zd/internal/scheduler"
	"zd/internal/service"
	"zd/internal/utils"

	"github.com/redis/go-redis/v9"
)

func init() {
	// Load the environment variables
	utils.LoadEnvVars()
}

func main() {
	// Setting up Service Dependencies =
	// Create and configure the RabbitMQ instance ==
	rabbitMQ := rabbitmq.New()
	rmqConnectionString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		utils.Env.RMQ_USER,
		utils.Env.RMQ_PASS,
		utils.Env.RMQ_DOMAIN,
		utils.Env.RMQ_PORT,
	)
	err := rabbitMQ.Connect(rmqConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	err = rabbitMQ.DeclareExchange("zendesk", "topic")
	if err != nil {
		log.Fatalf("Failed to decalre an exchange: %v", err)
	}

	// Creates the 2 routes that gets events from the exchange ==
	userEventIDRoute := rabbitMQ.RegisterExchangeRoute("userevent", rabbitmq.RouteTypeUserEventIDData)
	userEventNameRoute := rabbitMQ.RegisterExchangeRoute("marquee", rabbitmq.RouteTypeUserEventNameData)

	// Configuring the Core logic =
	// Create and configure the Redis Cache for the Core ==
	redisDB, err := strconv.Atoi(utils.Env.REDIS_DB)
	if err != nil {
		log.Fatalf("Failed to get the Redis Database index: %q", err)
	}
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", utils.Env.REDIS_DOMAIN, utils.Env.REDIS_PORT),
		Password: utils.Env.REDIS_PASS,
		DB:       redisDB,
	}
	redisCache := rediscache.NewRedisCache(opts)

	// Create the Service Core =
	core := core.NewZendeskMock(
		redisCache,
		fmt.Sprintf("%s:%s", utils.Env.USER_SRV_DOMAIN, utils.Env.USER_SRV_PORT),
		"/api/v2/events",
		"/api/v2/users",
	)

	// Create and configure the Service =
	srv := service.New(
		rabbitMQ,
		core,
	)

	// Registers the routes as callback functions in the service ==
	srv.RegisterPublishingCallback(userEventIDRoute.Publish, service.CallbackTypeImmediate)
	srv.RegisterPublishingCallback(userEventNameRoute.Publish, service.CallbackTypeLatest)

	// Create a scheduler for each of the required publishing schedules ==
	userEventIDScheduler := scheduler.New(3, true, srv.PublishNewUserEvent)
	userEventNameScheduler := scheduler.New(10, false, srv.PublishLatestUserEvent)

	// Run the scheduled tasks =
	userEventIDScheduler.Run()
	userEventNameScheduler.Run()

	// Start the Graceful Shutdown Channel =
	utils.GracefulShutdown([]utils.Closable{
		rabbitMQ,
	})

	// Loop to prevent main routine from stopping =
	var forever chan struct{}
	<-forever
}

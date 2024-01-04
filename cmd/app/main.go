package main

import (
	"fmt"
	"log"
	"zd/envvars"
	"zd/internal/core/service/zendeskservice"
	"zd/internal/driven/batch"
	"zd/internal/driven/rabbitmq"
	"zd/internal/drivers/ginserver"
	"zd/internal/drivers/schedule"
	"zd/internal/utils"

	"github.com/gin-gonic/gin"
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

	batch := batch.New()

	// Creating the Core Domain Service
	srv := zendeskservice.New(
		rabbitMQ,
		batch,
		fmt.Sprintf("%s:%s", envvars.Env.USER_SRV_DOMAIN, envvars.Env.USER_SRV_PORT),
		"/api/v1/event",
		"/api/v1/user",
	)

	// Creating and Configuring the Driver Actors
	//   HTTP Driver
	httpserver := ginserver.New(srv)
	if envvars.Env.ENV == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()
	server.GET("/", httpserver.GetUserEvent)

	// 	 Scheduler Driver
	generateUserEventSchedule := schedule.New(srv, 3, true, srv.GenerateUserEvent)
	drainSch := schedule.New(srv, 10, false, srv.PublishBatch)

	// Run all of the Drivers
	go func() {
		err := server.Run(envvars.Env.HTTP_PORT)
		if err != nil {
			panic(err)
		}
	}()

	generateUserEventSchedule.Run()
	drainSch.Run()

	// Starting Graceful Shutdown Channel
	utils.GracefulShutdown([]utils.Closable{
		rabbitMQ,
	})

	// Loop to prevent main routine from stopping
	var forever chan struct{}
	<-forever
}

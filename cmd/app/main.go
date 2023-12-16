package main

import (
	"fmt"
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
	queue := rabbitmq.New()
	rmqConnectionString := fmt.Sprintf(
		"amqp://%s:%s@%s:%s",
		envvars.Env.RMQ_USER,
		envvars.Env.RMQ_PASS,
		envvars.Env.RMQ_DOMAIN,
		envvars.Env.RMQ_PORT,
	)
	queue.Connect(rmqConnectionString)
	queue.DeclareExchange("zendesk", "topic")
	batch := batch.New()

	// Creating the Core Domain Service
	srv := zendeskservice.New(
		queue,
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

	// 	 Schedule Driver
	// genSch := schedule.New(srv, 50, srv.GenerateUserEvent)
	batchSch := schedule.New(srv, 10, srv.BatchUserEvent)
	drainSch := schedule.New(srv, 50, srv.PublishBatch)

	// Run the Drivers
	go func() {
		err := server.Run(envvars.Env.HTTP_PORT)
		if err != nil {
			panic(err)
		}
	}()

	// genSch.Run()
	batchSch.Run()
	drainSch.Run()

	// Starting Graceful Shutdown Channel
	utils.GracefulShutdown([]utils.Closable{
		queue,
	})

	// Loop to prevent main routine from stopping
	var forever chan struct{}
	<-forever
}

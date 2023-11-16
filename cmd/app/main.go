package main

import (
	"zd/envvars"
	"zd/internal/core/service/zendeskservice"
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
	queue := rabbitmq.New()
	queue.Connect("amqp://admin:admin@broker:5672")
	queue.DeclareExchange("zendesk", "topic")

	srv := zendeskservice.New(queue)
	httpserver := ginserver.NewGinServer(srv)
	scheduler := schedule.NewSchedule(srv, 50)

	if envvars.Env.ENV == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()
	server.GET("/", httpserver.GetUserEvent)

	utils.GracefulShutdown([]utils.Closable{
		queue,
	})

	go func() {
		err := server.Run(envvars.Env.PORT)
		if err != nil {
			panic(err)
		}
	}()

	scheduler.Run()

	var forever chan struct{}
	<-forever
}

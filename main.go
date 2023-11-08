package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aldrickdev/gin_api_template/envvars"
	"github.com/gin-gonic/gin"
)

func init() {
	envvars.LoadEnvVars()
}

func main() {
	fmt.Printf("Hello Gin, running in %s environment\n", envvars.Env.ENV)

	if envvars.Env.ENV == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/ping", Ping)

	err := r.Run(envvars.Env.PORT)
	if err != nil {
		panic(err)
	}
}

func Ping(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("Body: %s, Error: %v", string(body), err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

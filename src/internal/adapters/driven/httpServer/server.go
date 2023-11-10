package httpserver

import (
	"zd/envvars"
	"zd/internal/ports"

	"github.com/gin-gonic/gin"
)

type Adapter struct {
	api ports.APIPort
}

func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

func (s Adapter) Run() {
	if envvars.Env.ENV == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()
	registerRoutes(server, s)
	err := server.Run(envvars.Env.PORT)
	if err != nil {
		panic(err)
	}
}

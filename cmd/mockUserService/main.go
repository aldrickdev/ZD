package main

import (
	"fmt"
	"zd/internal/utils"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

var (
	MockDB = struct {
		Users  []User
		Events []Event
		Pods   []Pod
	}{
		Users: []User{
			{
				Id:    1,
				Name:  "Daniel",
				PodId: 1,
				Pod:   Containers,
			},
			{
				Id:    2,
				Name:  "Aldrick",
				PodId: 2,
				Pod:   Monitors,
			},
			{
				Id:    3,
				Name:  "Juliana",
				PodId: 2,
				Pod:   Monitors,
			},
			{
				Id:    4,
				Name:  "Danny",
				PodId: 1,
				Pod:   Containers,
			},
			{
				Id:    5,
				Name:  "Patrick",
				PodId: 1,
				Pod:   Containers,
			},
			{
				Id:    6,
				Name:  "Kevin",
				PodId: 3,
				Pod:   Cloud,
			},
		},
		Events: []Event{
			{
				Id:     1,
				Name:   "flare",
				Points: 10,
			},
			{
				Id:     2,
				Name:   "call",
				Points: 15,
			},
			{
				Id:     3,
				Name:   "screenshot",
				Points: 2,
			},
			{
				Id:     4,
				Name:   "datadog.yaml",
				Points: 1,
			},
			{
				Id:     5,
				Name:   "values.yaml",
				Points: 1,
			},
		},
		Pods: []Pod{
			Containers,
			Monitors,
			Cloud,
			Agent,
			WebPlatform,
			APM,
		},
	}
)

type User struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	PodId uint   `json:"pod_id"`
	Pod   Pod    `json:"pod"`
}
type Pod struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Points uint   `json:"points"`
}

var Containers = Pod{
	Id:   1,
	Name: "Containers",
}
var Monitors = Pod{
	Id:   2,
	Name: "Monitors",
}
var Cloud = Pod{
	Id:   3,
	Name: "Cloud",
}
var Agent = Pod{
	Id:   4,
	Name: "Agent",
}
var WebPlatform = Pod{
	Id:   5,
	Name: "WebPlatform",
}
var APM = Pod{
	Id:   6,
	Name: "APM",
}

type Event struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Points uint   `json:"points"`
}

func init() {
	utils.LoadEnvVars()
}

func main() {
	apiServer := echo.New()
	apiServer.Use(middleware.Logger())
	apiServer.Logger.SetLevel(log.DEBUG)
	apiServer.HideBanner = true

	apiServer.GET("/api/v2/users", func(c echo.Context) error {
		return c.JSON(200, MockDB.Users)
	})

	apiServer.GET("/api/v2/events", func(c echo.Context) error {
		return c.JSON(200, MockDB.Events)
	})
	apiServer.Logger.Info("Starting here")
	apiServer.Logger.Fatal(apiServer.Start(fmt.Sprintf(":%s", utils.Env.USER_SRV_PORT)))
}

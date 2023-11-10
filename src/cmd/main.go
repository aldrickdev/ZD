package main

import (
	"zd/envvars"
	httpserver "zd/internal/adapters/driven/httpServer"
	"zd/internal/applications/core/zendesk"
)

func init() {
	envvars.LoadEnvVars()
}

func main() {
	ep := zendesk.NewZendeskMock()
	httpServer := httpserver.NewAdapter(*ep)

	httpServer.Run()
}

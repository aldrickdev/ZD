package utils

import (
	"os"
)

type envVars struct {
	ENV             string
	USER_SRV_DOMAIN string
	USER_SRV_PORT   string
	RMQ_USER        string
	RMQ_PASS        string
	RMQ_DOMAIN      string
	RMQ_PORT        string
	REDIS_DOMAIN    string
	REDIS_PORT      string
	REDIS_DB        string
	REDIS_PASS      string
}

var Env envVars

func LoadEnvVars() {
	Env.ENV = os.Getenv("ENV")
	Env.USER_SRV_DOMAIN = os.Getenv("USER_SRV_DOMAIN")
	Env.USER_SRV_PORT = os.Getenv("USER_SRV_PORT")
	Env.RMQ_USER = os.Getenv("RMQ_USER")
	Env.RMQ_PASS = os.Getenv("RMQ_PASS")
	Env.RMQ_DOMAIN = os.Getenv("RMQ_DOMAIN")
	Env.RMQ_PORT = os.Getenv("RMQ_PORT")
	Env.REDIS_DOMAIN = os.Getenv("REDIS_DOMAIN")
	Env.REDIS_PORT = os.Getenv("REDIS_PORT")
	Env.REDIS_DB = os.Getenv("REDIS_DB")
	Env.REDIS_PASS = os.Getenv("REDIS_PASS")
}

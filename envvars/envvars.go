package envvars

import (
	"fmt"
	"os"
)

type envVars struct {
	ENV  string
	PORT string
}

var Env envVars

func LoadEnvVars() {
	Env.ENV = os.Getenv("ENV")
	Env.PORT = fmt.Sprintf(":%s", os.Getenv("PORT"))
}

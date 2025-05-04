package config

import (
	"github.com/Netflix/go-env"
)

type EnvVariables struct {
	Port     string `env:"PORT,default=3000"`
	GrpcPort string `env:"GRPC_PORT,default=50051"`
}

func GetEnvVariables() (*EnvVariables, error) {
	var envVars EnvVariables
	if _, err := env.UnmarshalFromEnviron(&envVars); err != nil {
		return nil, err
	}
	return &envVars, nil
}

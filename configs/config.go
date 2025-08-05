package configs

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Address    string `env:"ADDRESS,required"`
	PrivateKey string `env:"PRIVATE_KEY,required"`
	RpcUrl     string `env:"RPC_URL,required"`
}

func NewEnvConfig() *EnvConfig {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	config := &EnvConfig{}

	if err := env.Parse(config); err != nil {
		panic(err)
	}

	return config
}

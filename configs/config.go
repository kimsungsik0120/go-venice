package configs

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port       string `env:"PORT,required"`
	ChainId    int64  `env:"CHAIN_ID,required"`
	Address    string `env:"ADDRESS,required"`
	PrivateKey string `env:"PRIVATE_KEY,required"`
	RpcUrl     string `env:"RPC_URL,required"`
}

func Load() *EnvConfig {
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

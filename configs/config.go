package configs

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"os"
	"slices"
	"sync"
)

var (
	once sync.Once
	cfg  *EnvConfig
)

type EnvConfig struct {
	Port       string `env:"PORT,required"`
	ChainId    int64  `env:"CHAIN_ID,required"`
	Address    string `env:"ADDRESS,required"`
	PrivateKey string `env:"PRIVATE_KEY,required"`
	RpcUrl     string `env:"RPC_URL,required"`
}

func Load() *EnvConfig {
	once.Do(func() {
		goEnv := getenv("ENV", "local")
		if !slices.Contains([]string{"local", "dev", "prod"}, goEnv) {
			log.Fatalf("invalid ENV: %s", goEnv)
		}

		err := godotenv.Load("./env/." + goEnv + ".env")

		if err != nil {
			panic(err)
		}
		cfg = &EnvConfig{}
		if err := env.Parse(cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

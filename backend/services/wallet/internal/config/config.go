package config

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"

	sdkpostgre "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
)

// Config holds configuration for the project.
type Config struct {
	ServiceName              string `env:"SERVICE_NAME,default=wallet"`
	Env                      string `env:"ENV,default=development"`
	Port                     string `env:"PORT,default=9002"`
	Supabase                 Supabase
	Postgre                  sdkpostgre.Config
	GlobalTimeoutInSeconds   int `env:"GLOBAL_TIMEOUT_IN_SECONDS,default=60"`
	GracefulTimeoutInSeconds int `env:"GRACEFUL_TIMEOUT_IN_SECONDS,default=5"`
}

// Supabase holds config for Supabase.
type Supabase struct {
	JwksURL string `env:"SUPABASE_JWKS_URL,required"`
}

// NewConfig creates an instance of Config.
func NewConfig(ctx context.Context, lookuper envconfig.Lookuper, env string) *Config {
	_ = godotenv.Load(env)

	if lookuper == nil {
		lookuper = envconfig.OsLookuper()
	}

	var cfg Config
	err := envconfig.ProcessWith(ctx, &envconfig.Config{
		Target:   &cfg,
		Lookuper: lookuper,
	})
	if err != nil {
		panic(err)
	}
	return &cfg
}

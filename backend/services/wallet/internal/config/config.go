package config

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"

	sdkpostgre "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
	sdkmetric "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/metric"
	sdktrace "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/trace"
)

// Config holds configuration for the project.
type Config struct {
	Tracer                   sdktrace.Config
	Metric                   sdkmetric.Config
	ServiceName              string `env:"SERVICE_NAME,default=wallet"`
	Env                      string `env:"ENV,default=development"`
	Port                     string `env:"PORT,default=9002"`
	TopupSuccessURL          string `env:"TOPUP_SUCCESS_URL,default=http://localhost:9000"`
	Supabase                 Supabase
	Stripe                   Stripe
	Postgre                  sdkpostgre.Config
	GlobalTimeoutInSeconds   int `env:"GLOBAL_TIMEOUT_IN_SECONDS,default=60"`
	GracefulTimeoutInSeconds int `env:"GRACEFUL_TIMEOUT_IN_SECONDS,default=5"`
}

// Supabase holds config for Supabase.
type Supabase struct {
	JwksURL string `env:"SUPABASE_JWKS_URL,required"`
}

// Stripe holds config for Stripe.
type Stripe struct {
	APIKey string `env:"STRIPE_API_KEY,required"`
}

// New creates an instance of Config.
func New(ctx context.Context, lookuper envconfig.Lookuper, env string) *Config {
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

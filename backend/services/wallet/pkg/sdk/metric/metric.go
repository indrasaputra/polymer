package metric

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	interval = 5 * time.Second
)

// Config holds configuration for tracing.
type Config struct {
	OtelCollectorAddress string `env:"OPENTELEMETRY_COLLECTOR_ADDRESS,default=localhost:4317"`
	AppEnv               string `env:"APP_ENV,default=development"`
	ServiceName          string `env:"SERVICE_NAME,default=wallet"`
}

// Provider provides tracing functionality.
type Provider struct {
	*sdkmetric.MeterProvider
}

// NewProvider creates an instance of Provider.
func NewProvider(ctx context.Context, cfg Config) (*Provider, error) {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(cfg.OtelCollectorAddress),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.DeploymentEnvironmentKey.String(cfg.AppEnv),
		),
	)
	if err != nil {
		return nil, err
	}

	reader := sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(interval))
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return &Provider{provider}, nil
}

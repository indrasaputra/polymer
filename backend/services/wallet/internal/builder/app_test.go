package builder_test

import (
	"testing"

	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/internal/builder"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
	mockuow "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/pkg/sdk/uow"
)

func TestBuildWalletController(t *testing.T) {
	t.Run("success create wallet controller", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{},
		}

		handler := builder.BuildWalletController(dep)

		assert.NotNil(t, handler)
	})
}

func TestBuildWebhookController(t *testing.T) {
	t.Run("fail create webhook controller due to insufficient config", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{
				Kafka:  config.Kafka{Brokers: []string{"localhost:9092"}},
				Stripe: config.Stripe{},
			},
		}

		handler, err := builder.BuildWebhookController(dep)

		assert.Error(t, err)
		assert.Nil(t, handler)
	})
}

func TestBuildStripeEventConsumer(t *testing.T) {
	t.Run("fail create stripe event consumer due to insufficient config", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{
				Kafka:  config.Kafka{Brokers: []string{"localhost:9092"}},
				Stripe: config.Stripe{},
			},
		}

		consumer, err := builder.BuildStripeEventConsumer(dep)

		assert.Error(t, err)
		assert.Nil(t, consumer)
	})

	t.Run("success create stripe event consumer", func(t *testing.T) {
		dep := &builder.Dependency{
			Config: &config.Config{
				Kafka: config.Kafka{
					Brokers:                      []string{"localhost:9092"},
					StripeWebhookTopic:           "stripe-webhook",
					StripeWebhookConsumerGroupID: "wallet-consumer-group",
				},
				Stripe: config.Stripe{
					APIKey:        "key",
					WebhookSecret: "secret",
				},
			},
			StripeClient: builder.BuildStripeClient(&config.Config{
				Stripe: config.Stripe{
					APIKey:        "key",
					WebhookSecret: "secret",
				},
			}),
		}

		consumer, err := builder.BuildStripeEventConsumer(dep)

		assert.NoError(t, err)
		assert.NotNil(t, consumer)
	})
}

func TestBuildQueries(t *testing.T) {
	t.Run("success create queries", func(t *testing.T) {
		pool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %v\n", err)
		}
		g := mockuow.NewMockTxGetter(t)

		queries := builder.BuildQueries(pool, g)

		assert.NotNil(t, queries)
	})
}

func TestBuildPostgreWallet(t *testing.T) {
	t.Run("success create postgre wallet", func(t *testing.T) {
		pool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatalf("error opening a stub database connection: %v\n", err)
		}
		g := mockuow.NewMockTxGetter(t)
		queries := builder.BuildQueries(pool, g)

		wallet := builder.BuildPostgreWallet(queries)

		assert.NotNil(t, wallet)
	})
}

func TestBuildStripeClient(t *testing.T) {
	t.Run("success create stripe client", func(t *testing.T) {
		cfg := &config.Config{
			Stripe: config.Stripe{},
		}

		client := builder.BuildStripeClient(cfg)

		assert.NotNil(t, client)
	})
}

func TestBuildKafkaClient(t *testing.T) {
	t.Run("error create kafka client", func(t *testing.T) {
		cfg := &config.Config{
			Kafka: config.Kafka{},
		}

		client, err := builder.BuildKafkaClient(cfg)

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}

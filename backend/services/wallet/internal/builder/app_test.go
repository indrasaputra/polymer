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

func TestBuildStripeClient(t *testing.T) {
	t.Run("error create stripe client", func(t *testing.T) {
		cfg := &config.Config{
			Stripe: config.Stripe{},
		}

		client := builder.BuildStripeClient(cfg)

		assert.NotNil(t, client)
	})
}

package metric_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/metric"
)

func TestNewProvider(t *testing.T) {
	ctx := context.Background()

	t.Run("success create a new provider", func(t *testing.T) {
		cfg := metric.Config{
			OtelCollectorAddress: "localhost:4317",
			AppEnv:               "production",
		}

		prov, err := metric.NewProvider(ctx, cfg)

		assert.NoError(t, err)
		assert.NotNil(t, prov)
	})
}

package config_test

import (
	"context"
	"testing"

	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
)

var (
	testCtx = context.Background()
)

func TestNew(t *testing.T) {
	t.Run("fail load config due to missing required env", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = config.New(testCtx, nil, "")
		})
	})

	t.Run("success load config via .env", func(t *testing.T) {
		cfg := config.New(testCtx, nil, "../../env.example")

		assert.NotNil(t, cfg)
	})

	t.Run("success load config via lookupper", func(t *testing.T) {
		lookuper := envconfig.MapLookuper(map[string]string{
			"SUPABASE_JWKS_URL": "url",
			"POSTGRE_URL":       "url",
		})

		cfg := config.New(testCtx, lookuper, "")

		assert.NotNil(t, cfg)
	})
}

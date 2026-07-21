package builder

import (
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/client"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/controller"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/messaging"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/repository/db"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/repository/postgre"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	sdkpostgre "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/uow"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config       *config.Config
	TxManager    uow.TxManager
	Queries      *db.Queries
	StripeClient *client.Stripe
	KafkaClient  *messaging.Kafka
}

// BuildWalletController builds wallet controller including all of its dependencies.
func BuildWalletController(dep *Dependency) *controller.Wallet {
	pg := postgre.NewWallet(dep.Queries)
	creator := service.NewWalletCreator(dep.TxManager, pg, dep.StripeClient)
	topup := service.NewWalletTopup(pg, dep.StripeClient, dep.Config.TopupSuccessURL)
	return controller.NewWallet(creator, topup)
}

// BuildWebhookController builds webhook controller including all of its dependencies.
func BuildWebhookController(dep *Dependency) (*controller.Webhook, error) {
	cfg := service.StripeWebhookReceiverConfig{
		Producer:         dep.KafkaClient,
		EventConstructor: dep.StripeClient,
		Topic:            dep.Config.Stripe.WebhookTopic,
	}
	receiver, err := service.NewStripeWebhookReceiver(cfg)
	if err != nil {
		return nil, err
	}

	return controller.NewWebhook(receiver), nil
}

// BuildQueries builds sqlc queries.
func BuildQueries(tr uow.Tr, getter uow.TxGetter) *db.Queries {
	tx := sdkpostgre.NewTxDB(tr, getter)
	return db.New(tx)
}

// BuildStripeClient builds Stripe client.
func BuildStripeClient(cfg *config.Config) *client.Stripe {
	return client.NewStripe(cfg.Stripe.APIKey, cfg.Stripe.WebhookSecret)
}

// BuildKafkaClient builds Kafka client.
func BuildKafkaClient(cfg *config.Config) (*messaging.Kafka, error) {
	c, err := messaging.NewKafka(cfg.Kafka.Brokers)
	if err != nil {
		return nil, err
	}
	return c, nil
}

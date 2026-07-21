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
	KafkaClient  *messaging.KafkaProducer
	PgWallet     *postgre.Wallet
}

// BuildWalletController builds wallet controller including all of its dependencies.
func BuildWalletController(dep *Dependency) *controller.Wallet {
	creator := service.NewWalletCreator(dep.TxManager, dep.PgWallet, dep.StripeClient)
	topup := service.NewWalletTopup(dep.PgWallet, dep.StripeClient, dep.Config.TopupSuccessURL)
	return controller.NewWallet(creator, topup)
}

// BuildWebhookController builds webhook controller including all of its dependencies.
func BuildWebhookController(dep *Dependency) (*controller.Webhook, error) {
	e := service.NewStripeEventHandler(dep.PgWallet)
	cfg := service.StripeWebhookHandlerConfig{
		Producer:         dep.KafkaClient,
		EventConstructor: dep.StripeClient,
		Topic:            dep.Config.Kafka.StripeWebhookTopic,
		EventHandler:     e,
	}
	handler, err := service.NewStripeWebhookHandler(cfg)
	if err != nil {
		return nil, err
	}

	return controller.NewWebhook(handler), nil
}

// BuildStripeEventConsumer builds Stripe event consumer.
func BuildStripeEventConsumer(dep *Dependency) (*messaging.KafkaStripeWebhookConsumer, error) {
	e := service.NewStripeEventHandler(dep.PgWallet)
	cfg := service.StripeWebhookHandlerConfig{
		Producer:         dep.KafkaClient,
		EventConstructor: dep.StripeClient,
		Topic:            dep.Config.Kafka.StripeWebhookTopic,
		EventHandler:     e,
	}
	h, err := service.NewStripeWebhookHandler(cfg)
	if err != nil {
		return nil, err
	}

	c, err := messaging.NewKafkaStripeWebhookConsumer(h, dep.Config.Kafka.Brokers, dep.Config.Kafka.StripeWebhookTopic, dep.Config.Kafka.StripeWebhookConsumerGroupID)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// BuildQueries builds sqlc queries.
func BuildQueries(tr uow.Tr, getter uow.TxGetter) *db.Queries {
	tx := sdkpostgre.NewTxDB(tr, getter)
	return db.New(tx)
}

// BuildPostgreWallet builds postgre repository.
func BuildPostgreWallet(q *db.Queries) *postgre.Wallet {
	return postgre.NewWallet(q)
}

// BuildStripeClient builds Stripe client.
func BuildStripeClient(cfg *config.Config) *client.Stripe {
	return client.NewStripe(cfg.Stripe.APIKey, cfg.Stripe.WebhookSecret)
}

// BuildKafkaClient builds Kafka client.
func BuildKafkaClient(cfg *config.Config) (*messaging.KafkaProducer, error) {
	c, err := messaging.NewKafkaProducer(cfg.Kafka.Brokers)
	if err != nil {
		return nil, err
	}
	return c, nil
}

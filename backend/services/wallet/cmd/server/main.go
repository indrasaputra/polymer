// Server main program
package main

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/indrasaputra/polymer/backend/services/wallet/internal/builder"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/router"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/server"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/messaging"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
	wmid "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/http/middleware"
	sdklog "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/log"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/metric"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/trace"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/uow"
)

func main() {
	ctx := context.Background()
	cfg := config.New(ctx, nil, ".env")

	logger := sdklog.NewSlogLogger(cfg.ServiceName)
	slog.SetDefault(logger)

	traceProvider, err := trace.NewProvider(ctx, cfg.Tracer)
	raiseErrorIfAny(err)

	metricProvider, err := metric.NewProvider(ctx, cfg.Metric)
	raiseErrorIfAny(err)

	pool, err := postgre.NewPgxPool(cfg.Postgre)
	raiseErrorIfAny(err)
	defer pool.Close()

	txm, err := uow.NewTxManager(pool)
	raiseErrorIfAny(err)

	queries := builder.BuildQueries(pool, uow.NewTxGetter())
	pgWallet := builder.BuildPostgreWallet(queries)

	stripeClient := builder.BuildStripeClient(cfg)
	kafkaClient, err := builder.BuildKafkaClient(cfg)
	raiseErrorIfAny(err)

	dep := &builder.Dependency{
		Config:       cfg,
		TxManager:    txm,
		Queries:      queries,
		StripeClient: stripeClient,
		KafkaClient:  kafkaClient,
		PgWallet:     pgWallet,
	}

	stripeConsumer, err := builder.BuildStripeEventConsumer(dep)
	raiseErrorIfAny(err)

	srv, err := server.New(cfg, logger, traceProvider, metricProvider)
	raiseErrorIfAny(err)

	registerRouterForAPIV1(srv, dep)

	ctx, stop := srv.PrepareForGracefulStop()
	defer func() {
		_ = traceProvider.Shutdown(ctx)
		_ = metricProvider.Shutdown(ctx)
		kafkaClient.Close()
		stripeConsumer.Close()
		stop()

		slog.Info("done shutting down")
	}()

	err = runAll(ctx, stop,
		func(ctx context.Context) error { return runStripeWebhookConsumer(ctx, stripeConsumer) },
		func(ctx context.Context) error { return runServer(ctx, srv, cfg) },
	)
	raiseErrorIfAny(err)
}

func runServer(ctx context.Context, srv *server.Server, cfg *config.Config) error {
	slog.Info("starting http server")
	return srv.StartWithGracefulStop(ctx, cfg)
}

func runStripeWebhookConsumer(ctx context.Context, consumer *messaging.KafkaStripeWebhookConsumer) error {
	slog.Info("starting Stripe webhook consumer")
	consumer.Consume(ctx)
	return nil
}

func runAll(ctx context.Context, stop context.CancelFunc, fns ...func(context.Context) error) error {
	var wg sync.WaitGroup
	errs := make([]error, len(fns))

	for i, fn := range fns {
		wg.Add(1)
		go func(i int, fn func(context.Context) error) {
			defer wg.Done()
			if err := fn(ctx); err != nil && !errors.Is(err, context.Canceled) {
				errs[i] = err
				stop()
			}
		}(i, fn)
	}

	wg.Wait()
	return errors.Join(errs...)
}

func registerRouterForAPIV1(srv *server.Server, dep *builder.Dependency) {
	walletController := builder.BuildWalletController(dep)
	webhookController, err := builder.BuildWebhookController(dep)
	raiseErrorIfAny(err)

	jwtmid, err := wmid.NewJwtMiddleware(dep.Config)
	raiseErrorIfAny(err)

	router.RegisterAPIV1(srv.Echo, jwtmid, walletController, webhookController)
}

func raiseErrorIfAny(err error) {
	if err != nil {
		panic(err)
	}
}

// Server main program
package main

import (
	"context"

	"github.com/indrasaputra/polymer/backend/services/wallet/internal/builder"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/router"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/server"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/uow"
)

func main() {
	ctx := context.Background()
	cfg := config.New(ctx, nil, ".env")

	pool, err := postgre.NewPgxPool(cfg.Postgre)
	raiseErrorIfAny(err)
	defer pool.Close()

	txm, err := uow.NewTxManager(pool)
	raiseErrorIfAny(err)

	queries := builder.BuildQueries(pool, uow.NewTxGetter())

	dep := &builder.Dependency{
		Config:    cfg,
		TxManager: txm,
		Queries:   queries,
	}

	srv, err := server.New(cfg)
	raiseErrorIfAny(err)

	registerRouterForAPIV1(srv, dep)

	ctx, stop := srv.PrepareForGracefulStop()
	defer stop()

	err = srv.StartWithGracefulStop(ctx, cfg)
	raiseErrorIfAny(err)
}

func registerRouterForAPIV1(srv *server.Server, dep *builder.Dependency) {
	walletController := builder.BuildWalletController(dep)

	router.RegisterAPIV1(srv.Echo, walletController)
}

func raiseErrorIfAny(err error) {
	if err != nil {
		panic(err)
	}
}

package builder

import (
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/config"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/http/controller"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/repository/db"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/repository/postgre"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/service"
	sdkpostgre "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/uow"
)

// Dependency holds any dependency to build full use cases.
type Dependency struct {
	Config    *config.Config
	TxManager uow.TxManager
	Queries   *db.Queries
}

// BuildWalletController builds wallet controller including all of its dependencies.
func BuildWalletController(dep *Dependency) *controller.Wallet {
	p := postgre.NewWallet(dep.Queries)
	c := service.NewWalletCreator(p)
	return controller.NewWallet(c)
}

// BuildQueries builds sqlc queries.
func BuildQueries(tr uow.Tr, getter uow.TxGetter) *db.Queries {
	tx := sdkpostgre.NewTxDB(tr, getter)
	return db.New(tx)
}

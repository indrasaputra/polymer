package postgre_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/entity"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/repository/db"
	"github.com/indrasaputra/polymer/backend/services/wallet/internal/repository/postgre"
	sdkpostgre "github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
	mockuow "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/pkg/sdk/uow"
)

var (
	testCtx = context.Background()
)

type WalletSuite struct {
	pgWallet *postgre.Wallet
	db       pgxmock.PgxPoolIface
	getter   *mockuow.MockTxGetter
}

func TestNewWallet(t *testing.T) {
	t.Run("successfully create an instance of Wallet", func(t *testing.T) {
		st := createWalletSuite(t)
		assert.NotNil(t, st.pgWallet)
	})
}

func TestWallet_InsertWallet(t *testing.T) {
	queryInsert := `INSERT INTO wallets \(id, user_id, balance, currency, created_at, updated_at, created_by, updated_by\)
				VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)
				ON CONFLICT \(user_id, currency\) WHERE deleted_at IS NULL DO NOTHING
				RETURNING id, user_id, balance, currency, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by`
	querySelect := `SELECT id, user_id, balance, currency, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM wallets
					WHERE user_id = \$1 AND currency = \$2 AND deleted_at IS NULL
					LIMIT 1`

	t.Run("nil wallets is prohibited", func(t *testing.T) {
		st := createWalletSuite(t)

		res, err := st.pgWallet.InsertWallet(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet, err)
		assert.Nil(t, res)
	})

	t.Run("wallet exists due to DO NOTHING, but somehow not found during select", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			WillReturnError(sdkpostgre.ErrNotFound)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.UserID, wallet.Currency).
			WillReturnError(sdkpostgre.ErrNotFound)

		res, err := st.pgWallet.InsertWallet(testCtx, wallet)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrEmptyWallet, err)
		assert.Nil(t, res)
	})

	t.Run("wallet exists due to DO NOTHING, but error during select", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			WillReturnError(sdkpostgre.ErrNotFound)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.UserID, wallet.Currency).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.InsertWallet(testCtx, wallet)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("wallet exists due to DO NOTHING and success when select", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			WillReturnError(sdkpostgre.ErrNotFound)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.UserID, wallet.Currency).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.DeletedAt, wallet.CreatedBy, wallet.UpdatedBy, wallet.DeletedBy))

		res, err := st.pgWallet.InsertWallet(testCtx, wallet)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("insert wallet returns error", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.InsertWallet(testCtx, wallet)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success insert wallet", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.CreatedBy, wallet.UpdatedBy).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.DeletedAt, wallet.CreatedBy, wallet.UpdatedBy, wallet.DeletedBy))

		res, err := st.pgWallet.InsertWallet(testCtx, wallet)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWallet_InsertCustomer(t *testing.T) {
	queryInsert := `INSERT INTO customers \(id, user_id, stripe_customer_id, created_at, updated_at, created_by, updated_by\)
				VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7\)
				ON CONFLICT \(user_id\) WHERE deleted_at IS NULL DO NOTHING
				RETURNING id, user_id, stripe_customer_id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by`
	querySelect := `SELECT id, user_id, stripe_customer_id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM customers
					WHERE user_id = \$1 AND deleted_at IS NULL
					LIMIT 1`

	t.Run("nil customer is prohibited", func(t *testing.T) {
		st := createWalletSuite(t)

		res, err := st.pgWallet.InsertCustomer(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilCustomer, err)
		assert.Nil(t, res)
	})

	t.Run("customer exists due to DO NOTHING, but somehow not found during select", func(t *testing.T) {
		customer := createTestCustomer()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(customer.ID, customer.UserID, customer.StripeCustomerID, customer.CreatedAt, customer.UpdatedAt, customer.CreatedBy, customer.UpdatedBy).
			WillReturnError(sdkpostgre.ErrNotFound)
		st.db.ExpectQuery(querySelect).
			WithArgs(customer.UserID).
			WillReturnError(sdkpostgre.ErrNotFound)

		res, err := st.pgWallet.InsertCustomer(testCtx, customer)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilCustomer, err)
		assert.Nil(t, res)
	})

	t.Run("customer exists due to DO NOTHING, but error during select", func(t *testing.T) {
		customer := createTestCustomer()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(customer.ID, customer.UserID, customer.StripeCustomerID, customer.CreatedAt, customer.UpdatedAt, customer.CreatedBy, customer.UpdatedBy).
			WillReturnError(sdkpostgre.ErrNotFound)
		st.db.ExpectQuery(querySelect).
			WithArgs(customer.UserID).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.InsertCustomer(testCtx, customer)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("customer exists due to DO NOTHING and success when select", func(t *testing.T) {
		customer := createTestCustomer()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(customer.ID, customer.UserID, customer.StripeCustomerID, customer.CreatedAt, customer.UpdatedAt, customer.CreatedBy, customer.UpdatedBy).
			WillReturnError(sdkpostgre.ErrNotFound)
		st.db.ExpectQuery(querySelect).
			WithArgs(customer.UserID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "stripe_customer_id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(customer.ID, customer.UserID, customer.StripeCustomerID, customer.CreatedAt, customer.UpdatedAt, customer.DeletedAt, customer.CreatedBy, customer.UpdatedBy, customer.DeletedBy))

		res, err := st.pgWallet.InsertCustomer(testCtx, customer)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})

	t.Run("insert customer returns error", func(t *testing.T) {
		customer := createTestCustomer()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(customer.ID, customer.UserID, customer.StripeCustomerID, customer.CreatedAt, customer.UpdatedAt, customer.CreatedBy, customer.UpdatedBy).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.InsertCustomer(testCtx, customer)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success insert customer", func(t *testing.T) {
		customer := createTestCustomer()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(customer.ID, customer.UserID, customer.StripeCustomerID, customer.CreatedAt, customer.UpdatedAt, customer.CreatedBy, customer.UpdatedBy).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "stripe_customer_id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(customer.ID, customer.UserID, customer.StripeCustomerID, customer.CreatedAt, customer.UpdatedAt, customer.DeletedAt, customer.CreatedBy, customer.UpdatedBy, customer.DeletedBy))

		res, err := st.pgWallet.InsertCustomer(testCtx, customer)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWallet_GetCustomerByUserID(t *testing.T) {
	querySelect := `SELECT id, user_id, stripe_customer_id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM customers
					WHERE user_id = \$1 AND deleted_at IS NULL
					LIMIT 1`

	t.Run("customer not found", func(t *testing.T) {
		customer := createTestCustomer()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(customer.UserID).
			WillReturnError(sdkpostgre.ErrNotFound)

		res, err := st.pgWallet.GetCustomerByUserID(testCtx, customer.UserID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilCustomer, err)
		assert.Nil(t, res)
	})

	t.Run("query returns error", func(t *testing.T) {
		customer := createTestCustomer()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(customer.UserID).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.GetCustomerByUserID(testCtx, customer.UserID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success get customer", func(t *testing.T) {
		customer := createTestCustomer()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(customer.UserID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "stripe_customer_id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(customer.ID, customer.UserID, customer.StripeCustomerID, customer.CreatedAt, customer.UpdatedAt, customer.DeletedAt, customer.CreatedBy, customer.UpdatedBy, customer.DeletedBy))

		res, err := st.pgWallet.GetCustomerByUserID(testCtx, customer.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func createTestWallet() *entity.Wallet {
	userID := uuid.Must(uuid.NewV7())
	return &entity.Wallet{
		ID:       uuid.Must(uuid.NewV7()),
		UserID:   userID,
		Balance:  decimal.Zero,
		Currency: "USD",
		Auditable: entity.Auditable{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}
}

func createTestCustomer() *entity.Customer {
	userID := uuid.Must(uuid.NewV7())
	stripeID := "cus_NffrFeUfNV2Hib"
	return &entity.Customer{
		ID:               uuid.Must(uuid.NewV7()),
		UserID:           userID,
		StripeCustomerID: stripeID,
		Auditable: entity.Auditable{
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			CreatedBy: userID,
			UpdatedBy: userID,
		},
	}
}

func createWalletSuite(t *testing.T) *WalletSuite {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %v\n", err)
	}
	t.Cleanup(func() {
		assert.NoError(t, pool.ExpectationsWereMet())
		pool.Close()
	})
	g := mockuow.NewMockTxGetter(t)
	tx := sdkpostgre.NewTxDB(pool, g)
	q := db.New(tx)
	w := postgre.NewWallet(q)
	return &WalletSuite{
		pgWallet: w,
		db:       pool,
		getter:   g,
	}
}

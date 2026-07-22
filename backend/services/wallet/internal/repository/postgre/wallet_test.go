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

func TestWallet_GetActiveCustomerByUserID(t *testing.T) {
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

		res, err := st.pgWallet.GetActiveCustomerByUserID(testCtx, customer.UserID)

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

		res, err := st.pgWallet.GetActiveCustomerByUserID(testCtx, customer.UserID)

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

		res, err := st.pgWallet.GetActiveCustomerByUserID(testCtx, customer.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWallet_GetActiveWalletByIDAndUserID(t *testing.T) {
	querySelect := `SELECT id, user_id, balance, currency, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM wallets
					WHERE id = \$1 AND user_id = \$2 AND deleted_at IS NULL
					LIMIT 1`

	t.Run("wallet not found", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.ID, wallet.UserID).
			WillReturnError(sdkpostgre.ErrNotFound)

		res, err := st.pgWallet.GetActiveWalletByIDAndUserID(testCtx, wallet.ID, wallet.UserID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilWallet, err)
		assert.Nil(t, res)
	})

	t.Run("query returns error", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.ID, wallet.UserID).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.GetActiveWalletByIDAndUserID(testCtx, wallet.ID, wallet.UserID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success get wallet", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.ID, wallet.UserID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.DeletedAt, wallet.CreatedBy, wallet.UpdatedBy, wallet.DeletedBy))

		res, err := st.pgWallet.GetActiveWalletByIDAndUserID(testCtx, wallet.ID, wallet.UserID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWallet_InsertTransaction(t *testing.T) {
	queryInsert := `INSERT INTO transactions \(id, user_id, type, status, idempotency_key, amount, currency, checkout_session_id, created_at, updated_at, created_by, updated_by\)
				VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12\)
				RETURNING *`

	t.Run("nil transaction is prohibited", func(t *testing.T) {
		st := createWalletSuite(t)

		res, err := st.pgWallet.InsertTransaction(testCtx, nil)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilTransaction, err)
		assert.Nil(t, res)
	})

	t.Run("insert transaction returns error", func(t *testing.T) {
		trx := createTestTransaction()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(trx.ID, trx.UserID, db.TransactionType(trx.Type), db.TransactionStatus(trx.Status), trx.IdempotencyKey, trx.Amount, trx.Currency, trx.CheckoutSessionID, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.InsertTransaction(testCtx, trx)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success insert transaction", func(t *testing.T) {
		trx := createTestTransaction()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryInsert).
			WithArgs(trx.ID, trx.UserID, db.TransactionType(trx.Type), db.TransactionStatus(trx.Status), trx.IdempotencyKey, trx.Amount, trx.Currency, trx.CheckoutSessionID, trx.CreatedAt, trx.UpdatedAt, trx.CreatedBy, trx.UpdatedBy).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "type", "idempotency_key", "status", "amount", "currency", "checkout_session_id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(trx.ID, trx.UserID, db.TransactionType(trx.Type), db.TransactionStatus(trx.Status), trx.IdempotencyKey, trx.Amount, trx.Currency, trx.CheckoutSessionID, trx.CreatedAt, trx.UpdatedAt, trx.DeletedAt, trx.CreatedBy, trx.UpdatedBy, trx.DeletedBy))

		res, err := st.pgWallet.InsertTransaction(testCtx, trx)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWallet_GetActivePendingTransactionByIdempotencyKey(t *testing.T) {
	querySelect := `SELECT id, user_id, type, status, idempotency_key, amount, currency, checkout_session_id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM transactions
					WHERE idempotency_key = \$1 AND status = 'pending' AND deleted_at IS NULL
					LIMIT 1`

	t.Run("transaction not found", func(t *testing.T) {
		trx := createTestTransaction()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(trx.IdempotencyKey).
			WillReturnError(sdkpostgre.ErrNotFound)

		res, err := st.pgWallet.GetActivePendingTransactionByIdempotencyKey(testCtx, trx.IdempotencyKey)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilTransaction, err)
		assert.Nil(t, res)
	})

	t.Run("query returns error", func(t *testing.T) {
		trx := createTestTransaction()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(trx.IdempotencyKey).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.GetActivePendingTransactionByIdempotencyKey(testCtx, trx.IdempotencyKey)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success get pending transaction", func(t *testing.T) {
		trx := createTestTransaction()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(trx.IdempotencyKey).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "type", "idempotency_key", "status", "amount", "currency", "checkout_session_id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(trx.ID, trx.UserID, db.TransactionType(trx.Type), db.TransactionStatus(trx.Status), trx.IdempotencyKey, trx.Amount, trx.Currency, trx.CheckoutSessionID, trx.CreatedAt, trx.UpdatedAt, trx.DeletedAt, trx.CreatedBy, trx.UpdatedBy, trx.DeletedBy))

		res, err := st.pgWallet.GetActivePendingTransactionByIdempotencyKey(testCtx, trx.IdempotencyKey)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWallet_UpdateActiveTransactionToCompletedByCheckoutSessionID(t *testing.T) {
	queryUpdate := `UPDATE transactions
					SET status = 'completed', updated_at = \$1, updated_by = \$2
					WHERE checkout_session_id = \$3 AND deleted_at IS NULL
					RETURNING id, user_id, type, status, idempotency_key, amount, currency, checkout_session_id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by`

	t.Run("transaction not found", func(t *testing.T) {
		sessionID := "cs_test_123"
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryUpdate).
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), &sessionID).
			WillReturnError(sdkpostgre.ErrNotFound)

		err := st.pgWallet.UpdateActiveTransactionToCompletedByCheckoutSessionID(testCtx, sessionID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilTransaction, err)
	})

	t.Run("query returns error", func(t *testing.T) {
		sessionID := "cs_test_123"
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryUpdate).
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), &sessionID).
			WillReturnError(assert.AnError)

		err := st.pgWallet.UpdateActiveTransactionToCompletedByCheckoutSessionID(testCtx, sessionID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
	})

	t.Run("success update transaction to completed", func(t *testing.T) {
		sessionID := "cs_test_123"
		trx := createTestTransaction()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryUpdate).
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), &sessionID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "type", "idempotency_key", "status", "amount", "currency", "checkout_session_id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(trx.ID, trx.UserID, db.TransactionType(trx.Type), db.TransactionStatus(trx.Status), trx.IdempotencyKey, trx.Amount, trx.Currency, trx.CheckoutSessionID, trx.CreatedAt, trx.UpdatedAt, trx.DeletedAt, trx.CreatedBy, trx.UpdatedBy, trx.DeletedBy))

		err := st.pgWallet.UpdateActiveTransactionToCompletedByCheckoutSessionID(testCtx, sessionID)

		assert.NoError(t, err)
	})
}

func TestWallet_GetActiveTransactionByCheckoutSessionIDForUpdate(t *testing.T) {
	querySelect := `SELECT id, user_id, type, status, idempotency_key, amount, currency, checkout_session_id, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by FROM transactions
					WHERE checkout_session_id = \$1 AND deleted_at IS NULL
					LIMIT 1 FOR NO KEY UPDATE`

	t.Run("transaction not found", func(t *testing.T) {
		sessionID := "cs_test_123"
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(&sessionID).
			WillReturnError(sdkpostgre.ErrNotFound)

		res, err := st.pgWallet.GetActiveTransactionByCheckoutSessionIDForUpdate(testCtx, sessionID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilTransaction, err)
		assert.Nil(t, res)
	})

	t.Run("query returns error", func(t *testing.T) {
		sessionID := "cs_test_123"
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(&sessionID).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.GetActiveTransactionByCheckoutSessionIDForUpdate(testCtx, sessionID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success get active transaction for update", func(t *testing.T) {
		sessionID := "cs_test_123"
		trx := createTestTransaction()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(&sessionID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "type", "idempotency_key", "status", "amount", "currency", "checkout_session_id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(trx.ID, trx.UserID, db.TransactionType(trx.Type), db.TransactionStatus(trx.Status), trx.IdempotencyKey, trx.Amount, trx.Currency, trx.CheckoutSessionID, trx.CreatedAt, trx.UpdatedAt, trx.DeletedAt, trx.CreatedBy, trx.UpdatedBy, trx.DeletedBy))

		res, err := st.pgWallet.GetActiveTransactionByCheckoutSessionIDForUpdate(testCtx, sessionID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWallet_GetActiveWalletByIDForUpdate(t *testing.T) {
	querySelect := `SELECT id, user_id, balance, currency, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by
					FROM wallets WHERE id = \$1 AND deleted_at IS NULL  LIMIT 1 FOR NO KEY UPDATE`

	t.Run("wallet not found", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.ID).
			WillReturnError(sdkpostgre.ErrNotFound)

		res, err := st.pgWallet.GetActiveWalletByIDForUpdate(testCtx, wallet.ID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrNilWallet, err)
		assert.Nil(t, res)
	})

	t.Run("query returns error", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.ID).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.GetActiveWalletByIDForUpdate(testCtx, wallet.ID)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success get active wallet for update", func(t *testing.T) {
		wallet := createTestWallet()
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(querySelect).
			WithArgs(wallet.ID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.DeletedAt, wallet.CreatedBy, wallet.UpdatedBy, wallet.DeletedBy))

		res, err := st.pgWallet.GetActiveWalletByIDForUpdate(testCtx, wallet.ID)

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestWallet_AddActiveWalletBalance(t *testing.T) {
	queryUpdate := `UPDATE wallets SET balance = balance \+ \$2 WHERE id = \$1 AND deleted_at IS NULL --noqa
					RETURNING id, user_id, balance, currency, created_at, updated_at, deleted_at, created_by, updated_by, deleted_by`

	t.Run("add wallet balance returns error", func(t *testing.T) {
		wallet := createTestWallet()
		amount := decimal.NewFromInt(50)
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryUpdate).
			WithArgs(wallet.ID, amount).
			WillReturnError(assert.AnError)

		res, err := st.pgWallet.AddActiveWalletBalance(testCtx, wallet.ID, amount)

		assert.Error(t, err)
		assert.Equal(t, entity.ErrInternal, err)
		assert.Nil(t, res)
	})

	t.Run("success add wallet balance", func(t *testing.T) {
		wallet := createTestWallet()
		amount := decimal.NewFromInt(50)
		st := createWalletSuite(t)
		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery(queryUpdate).
			WithArgs(wallet.ID, amount).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "balance", "currency", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by"}).
				AddRow(wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency, wallet.CreatedAt, wallet.UpdatedAt, wallet.DeletedAt, wallet.CreatedBy, wallet.UpdatedBy, wallet.DeletedBy))

		res, err := st.pgWallet.AddActiveWalletBalance(testCtx, wallet.ID, amount)

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

func createTestTransaction() *entity.Transaction {
	userID := uuid.Must(uuid.NewV7())
	sessionID := "cs_test_123"
	return &entity.Transaction{
		ID:                uuid.Must(uuid.NewV7()),
		UserID:            userID,
		Type:              entity.TransactionType("topup"),
		Status:            entity.TransactionStatus("pending"),
		IdempotencyKey:    uuid.Must(uuid.NewV7()),
		Amount:            decimal.NewFromInt(100),
		Currency:          "USD",
		CheckoutSessionID: &sessionID,
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

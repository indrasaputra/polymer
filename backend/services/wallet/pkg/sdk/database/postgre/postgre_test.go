package postgre_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v5"
	"github.com/stretchr/testify/assert"

	"github.com/indrasaputra/polymer/backend/services/wallet/pkg/sdk/database/postgre"
	mockuow "github.com/indrasaputra/polymer/backend/services/wallet/test/mock/pkg/sdk/uow"
)

var (
	testCtx = context.Background()
)

type TxDBSuite struct {
	db     pgxmock.PgxPoolIface
	getter *mockuow.MockTxGetter
	tx     *postgre.TxDB
}

func TestIsUniqueViolationError(t *testing.T) {
	t.Run("error is not unique violation", func(t *testing.T) {
		res := postgre.IsUniqueViolationError(assert.AnError)
		assert.False(t, res)
	})

	t.Run("nil error is not unique violation", func(t *testing.T) {
		res := postgre.IsUniqueViolationError(nil)
		assert.False(t, res)
	})

	t.Run("error is unique violation", func(t *testing.T) {
		res := postgre.IsUniqueViolationError(&pgconn.PgError{Code: "23505"})
		assert.True(t, res)
	})
}

func TestNewPgxPool(t *testing.T) {
	t.Run("success create pgx pool", func(t *testing.T) {
		pool, err := postgre.NewPgxPool(postgre.Config{})

		assert.NoError(t, err)
		assert.NotNil(t, pool)
	})
}

func TestNewTxDB(t *testing.T) {
	t.Run("success create txdb", func(t *testing.T) {
		st := createTxDBSuite(t)

		assert.NotNil(t, st.tx)
	})
}

func TestTxDB_Exec(t *testing.T) {
	t.Run("success execute exec", func(t *testing.T) {
		st := createTxDBSuite(t)

		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectExec("exec").WithArgs("arg1", "arg2").WillReturnResult(pgxmock.NewResult("exec", 1))

		res, err := st.tx.Exec(testCtx, "exec", "arg1", "arg2")

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestTxDB_Query(t *testing.T) {
	t.Run("success execute query", func(t *testing.T) {
		st := createTxDBSuite(t)

		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery("query").
			WithArgs("arg1", "arg2").
			WillReturnRows(pgxmock.
				NewRows([]string{"id"}).
				AddRow("id"))

		res, err := st.tx.Query(testCtx, "query", "arg1", "arg2")

		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestTxDB_QueryRow(t *testing.T) {
	t.Run("success execute query row", func(t *testing.T) {
		st := createTxDBSuite(t)

		st.getter.EXPECT().DefaultTrOrDB(testCtx, st.db).Return(st.db)
		st.db.ExpectQuery("query").
			WithArgs("arg1", "arg2").
			WillReturnRows(pgxmock.
				NewRows([]string{"id"}).
				AddRow("id"))

		res := st.tx.QueryRow(testCtx, "query", "arg1", "arg2")

		assert.NotNil(t, res)
	})
}

func createTxDBSuite(t *testing.T) *TxDBSuite {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("error opening a stub database connection: %v\n", err)
	}
	g := mockuow.NewMockTxGetter(t)
	tx := postgre.NewTxDB(pool, g)
	return &TxDBSuite{
		db:     pool,
		getter: g,
		tx:     tx,
	}
}

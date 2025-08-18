package transaction

import (
	"context"

	"github.com/en7ka/chat-server/internal/client/db"
	"github.com/en7ka/chat-server/internal/client/db/pg"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {

		// восстанавливаемся после паники
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		//откатываем транзакцию, если есть ошибка
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "rollback failed: %v", errRollback)
			}

			return
		}

		//если нет ошибок, то коммит
		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "commit failed")
			}
		}

	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "transaction failed")
	}

	return err
}

func (m *manager) ReadCommited(ctx context.Context, f db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, f)
}

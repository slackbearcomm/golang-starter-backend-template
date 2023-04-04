package dbstore

import (
	"context"
	"fmt"
	"gogql/utils/faulterr"
	"gogql/utils/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX struct {
	conn *pgxpool.Pool
}

func NewDBTX(conn *pgxpool.Pool) *DBTX {
	return &DBTX{conn}
}

func (t *DBTX) BeginTx(ctx context.Context) (pgx.Tx, *faulterr.FaultErr) {
	txOptions := pgx.TxOptions{
		IsoLevel:       pgx.TxIsoLevel(pgx.ReadCommitted),
		AccessMode:     pgx.TxAccessMode(pgx.ReadWrite),
		DeferrableMode: pgx.TxDeferrableMode(pgx.NotDeferrable),
	}

	tx, txErr := t.conn.BeginTx(ctx, txOptions)
	if txErr != nil {
		return nil, faulterr.NewInternalServerError(txErr.Error())
	}
	logger.Info("Beign Transaction")

	return tx, nil

}

func (t *DBTX) CommitTx(ctx context.Context, tx pgx.Tx) *faulterr.FaultErr {
	err := tx.Commit(ctx)
	if err != nil {
		return faulterr.NewInternalServerError(err.Error())
	}
	logger.Info("Commit Transaction")

	return nil
}

func (t *DBTX) RollbackTx(ctx context.Context, tx pgx.Tx) *faulterr.FaultErr {
	err := tx.Rollback(ctx)
	if err != nil {
		// return faulterr.NewInternalServerError(err.Error())
		logger.Info(fmt.Sprintf("Rollback Transaction %s", err.Error()))
		return nil
	}
	logger.Warning("Rollback Transaction")

	return nil
}

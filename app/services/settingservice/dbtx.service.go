package settingservice

import (
	"context"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/jackc/pgx/v5"
)

// This DBTX acts just like an interface between resolvers and dbstore
// as resolvers and dbstore should not be directly exposed
type DBTX struct {
	dbStore *dbstore.DBStore
}

func NewDBTX(dbStore *dbstore.DBStore) *DBTX {
	return &DBTX{dbStore}
}

func (s *DBTX) BeginTx(ctx context.Context) (pgx.Tx, *faulterr.FaultErr) {
	return s.dbStore.DBTX.BeginTx(ctx)

}

func (s *DBTX) CommitTx(ctx context.Context, tx pgx.Tx) *faulterr.FaultErr {
	return s.dbStore.DBTX.CommitTx(ctx, tx)
}

func (s *DBTX) RollbackTx(ctx context.Context, tx pgx.Tx) *faulterr.FaultErr {
	return s.dbStore.DBTX.RollbackTx(ctx, tx)
}

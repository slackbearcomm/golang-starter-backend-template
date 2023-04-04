package orgmaster

import (
	"context"
	"gogql/app/helpers"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"
	"time"

	"github.com/jackc/pgx/v5"
)

type AuthSessionMaster struct {
	dbstore *dbstore.DBStore
}

func NewAuthSessionMaster(s *dbstore.DBStore) *AuthSessionMaster {
	return &AuthSessionMaster{s}
}

func (m *AuthSessionMaster) Create(ctx context.Context, tx pgx.Tx, userID int64) (*dbmodels.AuthSession, *faulterr.FaultErr) {
	obj, err := m.construct(userID)
	if err != nil {
		return nil, err
	}
	return m.dbstore.AuthSessionStore.Insert(ctx, tx, obj)
}

func (m *AuthSessionMaster) construct(userID int64) (*dbmodels.AuthSession, *faulterr.FaultErr) {
	token, err := helpers.GenerateUID()
	if err != nil {
		return nil, err
	}

	obj := &dbmodels.AuthSession{
		UserID:    userID,
		Token:     *token,
		IsValid:   true,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(720)), // token will expire after 5 minutes
	}

	return obj, nil
}

package orgmaster

import (
	"context"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/encrypt"
	"gogql/utils/faulterr"
	"time"

	"github.com/jackc/pgx/v5"
)

type OTPSessionMaster struct {
	dbstore *dbstore.DBStore
}

func NewOTPSessionMaster(s *dbstore.DBStore) *OTPSessionMaster {
	return &OTPSessionMaster{s}
}

func (m *OTPSessionMaster) Create(ctx context.Context, tx pgx.Tx, userID int64) (*dbmodels.OTPSession, *faulterr.FaultErr) {
	return m.dbstore.OTPSessionStore.Insert(ctx, tx, m.construct(userID))
}

func (m *OTPSessionMaster) construct(userID int64) *dbmodels.OTPSession {
	return &dbmodels.OTPSession{
		UserID:    userID,
		Token:     encrypt.GenerateRandomString(5),
		IsValid:   true,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(5)), // token will expire after 5 minutes
	}
}

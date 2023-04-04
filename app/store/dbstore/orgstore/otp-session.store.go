package orgstore

import (
	"context"
	"gogql/app/models/dbmodels"
	"gogql/utils/faulterr"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OTPSessionStore struct {
	conn *pgxpool.Pool
}

var _ OTPSessionStoreInterface = &OTPSessionStore{}

type OTPSessionStoreInterface interface {
	GetByToken(ctx context.Context, token string) (*dbmodels.OTPSession, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, arg *dbmodels.OTPSession) (*dbmodels.OTPSession, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, arg *dbmodels.OTPSession) *faulterr.FaultErr
}

func NewOTPSessionStore(conn *pgxpool.Pool) *OTPSessionStore {
	return &OTPSessionStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetByToken gets otp session by code from database
func (s *OTPSessionStore) GetByToken(ctx context.Context, token string) (*dbmodels.OTPSession, *faulterr.FaultErr) {
	errMsg := "error when trying to get otp session by token"

	queryStmt := `
	SELECT * FROM otp_sessions
	WHERE otp_sessions.token = $1
	`

	row := s.conn.QueryRow(ctx, queryStmt, token)
	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert inserts a otp session in database
func (s *OTPSessionStore) Insert(ctx context.Context, tx pgx.Tx, arg *dbmodels.OTPSession) (*dbmodels.OTPSession, *faulterr.FaultErr) {
	errMsg := "error when trying to insert otp session"

	queryStmt := `
	INSERT INTO
	otp_sessions(
		user_id,
		token,
		is_valid,
		expires_at
	)
	VALUES ($1, $2, $3, $4)
	RETURNING *
	`

	row := tx.QueryRow(ctx, queryStmt,
		arg.UserID,
		arg.Token,
		arg.IsValid,
		arg.ExpiresAt,
	)

	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// Update User
func (s *OTPSessionStore) Update(ctx context.Context, tx pgx.Tx, arg *dbmodels.OTPSession) *faulterr.FaultErr {
	errMsg := "error when trying to update user"

	queryStmt := `
	UPDATE otp_sessions
	SET
		is_valid=$1
	WHERE id=$2
	`

	_, err := tx.Exec(ctx, queryStmt,
		&arg.IsValid,
		&arg.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}
	return nil
}

func (s *OTPSessionStore) scanRow(row pgx.Row) (*dbmodels.OTPSession, error) {
	obj := dbmodels.OTPSession{}

	if err := row.Scan(
		&obj.ID,
		&obj.UserID,
		&obj.Token,
		&obj.IsValid,
		&obj.ExpiresAt,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &obj, nil
}

package orgstore

import (
	"context"
	"gogql/app/models/dbmodels"
	"gogql/utils/faulterr"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthSessionStore struct {
	conn *pgxpool.Pool
}

var _ AuthSessionStoreInterface = &AuthSessionStore{}

type AuthSessionStoreInterface interface {
	GetByToken(ctx context.Context, token uuid.UUID) (*dbmodels.AuthSession, *faulterr.FaultErr)
	Insert(ctx context.Context, tx pgx.Tx, arg *dbmodels.AuthSession) (*dbmodels.AuthSession, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, arg *dbmodels.AuthSession) *faulterr.FaultErr
}

func NewAuthSessionStore(conn *pgxpool.Pool) *AuthSessionStore {
	return &AuthSessionStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetByToken gets auth session by code from database
func (s *AuthSessionStore) GetByToken(ctx context.Context, token uuid.UUID) (*dbmodels.AuthSession, *faulterr.FaultErr) {
	errMsg := "error when trying to get auth session by token"

	queryStmt := `
	SELECT * FROM auth_sessions
	WHERE auth_sessions.token = $1
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

// Insert inserts a auth session in database
func (s *AuthSessionStore) Insert(ctx context.Context, tx pgx.Tx, arg *dbmodels.AuthSession) (*dbmodels.AuthSession, *faulterr.FaultErr) {
	errMsg := "error when trying to insert auth session"

	queryStmt := `
	INSERT INTO
	auth_sessions(
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
func (s *AuthSessionStore) Update(ctx context.Context, tx pgx.Tx, arg *dbmodels.AuthSession) *faulterr.FaultErr {
	errMsg := "error when trying to update user"

	queryStmt := `
	UPDATE auth_sessions
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

func (s *AuthSessionStore) scanRow(row pgx.Row) (*dbmodels.AuthSession, error) {
	obj := dbmodels.AuthSession{}

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

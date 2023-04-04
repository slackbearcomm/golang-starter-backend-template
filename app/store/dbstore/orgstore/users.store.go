package orgstore

import (
	"context"
	"fmt"
	"gogql/app/models"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore/dbhelpers"
	"gogql/utils/faulterr"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStore struct {
	conn *pgxpool.Pool
}

var _ UserStoreInterface = &UserStore{}

type UserStoreInterface interface {
	GetManyByIDs(ctx context.Context, ids []int64) ([]*dbmodels.User, *faulterr.FaultErr)

	List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID, roleID *int64) ([]dbmodels.User, int, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*dbmodels.User, *faulterr.FaultErr)
	GetByEmail(ctx context.Context, email string) (*dbmodels.User, *faulterr.FaultErr)
	GetByPhone(ctx context.Context, phone string) (*dbmodels.User, *faulterr.FaultErr)

	Insert(ctx context.Context, tx pgx.Tx, u dbmodels.User) (*dbmodels.User, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, u dbmodels.User) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewUserStore(conn *pgxpool.Pool) *UserStore {
	return &UserStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetManyByIDs get all users by ids
func (s *UserStore) GetManyByIDs(ctx context.Context, ids []int64) ([]*dbmodels.User, *faulterr.FaultErr) {
	errMsg := "error when trying to get many users by ids"

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	queryStmt := "SELECT * from users WHERE id IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := s.conn.Query(ctx, queryStmt, args...)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	output, err := s.scanRows(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	result := []*dbmodels.User{}
	for i := 0; i < len(output); i++ {
		result = append(result, &output[i])
	}
	return result, nil
}

// List gets all users
func (s *UserStore) List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID, roleID *int64) ([]dbmodels.User, int, *faulterr.FaultErr) {
	errMsg := "error when trying to get users"

	// define query
	selectQuery := `SELECT * FROM users`
	conditionsQuery := `
	WHERE ($1::UUID IS NULL OR $1 = org_uid)
	AND ($2::INTEGER IS NULL OR $2 = role_id)
	AND ($3::BOOLEAN IS NULL OR $3 = is_final)
	AND ($4::BOOLEAN IS NULL OR $4 = is_archived)
	`
	filterQuery := dbhelpers.ResolveFilterSort(dbhelpers.UsersTable, filter)
	queryStmt := fmt.Sprintf("%s %s %s", selectQuery, conditionsQuery, filterQuery)

	var queryArgs []interface{}
	queryArgs = append(queryArgs, orgUID, roleID, filter.IsFinal, filter.IsArchived)

	// query rows
	rows, err := s.conn.Query(ctx, queryStmt, queryArgs...)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	// get total
	total, err := dbhelpers.GetGlobalTotal(ctx, s.conn, dbhelpers.UsersTable, conditionsQuery, queryArgs)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}

	result, err := s.scanRows(rows)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	return result, total, nil
}

// GetByID User
func (s *UserStore) GetByID(ctx context.Context, id int64) (*dbmodels.User, *faulterr.FaultErr) {
	errMsg := "error when trying to get user by id"

	queryStmt := `SELECT * FROM users WHERE id=$1`

	row := s.conn.QueryRow(ctx, queryStmt, id)
	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// GetByEmail User
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*dbmodels.User, *faulterr.FaultErr) {
	errMsg := fmt.Sprintf("error when trying to get user by email - %s", email)

	queryStmt := `SELECT * FROM users WHERE email=$1`

	row := s.conn.QueryRow(ctx, queryStmt, email)
	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// GetByPhone User
func (s *UserStore) GetByPhone(ctx context.Context, phone string) (*dbmodels.User, *faulterr.FaultErr) {
	errMsg := fmt.Sprintf("error when trying to get user by phone - %s", phone)

	queryStmt := `SELECT * FROM users WHERE phone=$1`

	row := s.conn.QueryRow(ctx, queryStmt, phone)
	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert User
func (s *UserStore) Insert(ctx context.Context, tx pgx.Tx, arg dbmodels.User) (*dbmodels.User, *faulterr.FaultErr) {
	errMsg := "error when trying to insert user"

	queryStmt := `
	INSERT INTO
	users(
		first_name,
		last_name,
		email,
		phone,
		is_admin,
		org_uid,
		role_id,
		status,
		is_final
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING *
	`

	row := tx.QueryRow(ctx, queryStmt,
		&arg.FirstName,
		&arg.LastName,
		&arg.Email,
		&arg.Phone,
		&arg.IsAdmin,
		&arg.OrgUID,
		&arg.RoleID,
		&arg.Status,
		&arg.IsFinal,
	)

	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// Update User
func (s *UserStore) Update(ctx context.Context, tx pgx.Tx, arg dbmodels.User) *faulterr.FaultErr {
	errMsg := "error when trying to update user"

	queryStmt := `
	UPDATE users
	SET
		first_name=$1,
		last_name=$2,
		email=$3,
		phone=$4,
		role_id=$5,
		status=$6,
		is_archived=$7
	WHERE id=$8
	`

	_, err := tx.Exec(ctx, queryStmt,
		&arg.FirstName,
		&arg.LastName,
		&arg.Email,
		&arg.Phone,
		&arg.RoleID,
		&arg.Status,
		&arg.IsArchived,
		&arg.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}
	return nil
}

// Delete User
func (s *UserStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM users WHERE id=$1`

	_, err := tx.Exec(ctx, queryStmt, id)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete user")
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *UserStore) scanRows(rows pgx.Rows) ([]dbmodels.User, error) {
	result := []dbmodels.User{}
	obj := dbmodels.User{}

	for rows.Next() {
		if err := rows.Scan(
			&obj.ID,
			&obj.FirstName,
			&obj.LastName,
			&obj.Email,
			&obj.Phone,
			&obj.IsAdmin,
			&obj.OrgUID,
			&obj.RoleID,
			&obj.Status,
			&obj.IsFinal,
			&obj.IsArchived,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, nil
}

func (s *UserStore) scanRow(row pgx.Row) (*dbmodels.User, error) {
	obj := &dbmodels.User{}
	if err := row.Scan(
		&obj.ID,
		&obj.FirstName,
		&obj.LastName,
		&obj.Email,
		&obj.Phone,
		&obj.IsAdmin,
		&obj.OrgUID,
		&obj.RoleID,
		&obj.Status,
		&obj.IsFinal,
		&obj.IsArchived,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return obj, nil
}

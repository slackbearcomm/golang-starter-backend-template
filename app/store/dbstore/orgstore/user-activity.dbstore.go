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

type UserActivityStore struct {
	conn *pgxpool.Pool
}

var _ UserActivityStoreInterface = &UserActivityStore{}

type UserActivityStoreInterface interface {
	GetManyByIDs(ctx context.Context, ids []int64) ([]*dbmodels.UserActivity, *faulterr.FaultErr)

	List(ctx context.Context, filter models.SearchFilter, userID *int64, orgUID *uuid.UUID) ([]dbmodels.UserActivity, int, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*dbmodels.UserActivity, *faulterr.FaultErr)

	Insert(ctx context.Context, tx pgx.Tx, u dbmodels.UserActivity) (*dbmodels.UserActivity, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, u dbmodels.UserActivity) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewUserActivityStore(conn *pgxpool.Pool) *UserActivityStore {
	return &UserActivityStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetManyByIDs get all user_activities by ids
func (s *UserActivityStore) GetManyByIDs(ctx context.Context, ids []int64) ([]*dbmodels.UserActivity, *faulterr.FaultErr) {
	errMsg := "error when trying to get many user_activities by ids"

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	queryStmt := "SELECT * from user_activities WHERE id IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := s.conn.Query(ctx, queryStmt, args...)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	output, err := s.scanRows(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	result := []*dbmodels.UserActivity{}
	for i := 0; i < len(output); i++ {
		result = append(result, &output[i])
	}
	return result, nil
}

// List gets all user_activities
func (s *UserActivityStore) List(ctx context.Context, filter models.SearchFilter, userID *int64, orgUID *uuid.UUID) ([]dbmodels.UserActivity, int, *faulterr.FaultErr) {
	errMsg := "error when trying to get user_activities"

	// define query
	selectQuery := `SELECT * FROM user_activities`
	conditionsQuery := `
	WHERE ($1::INTEGER IS NULL OR $1 = user_id)
	AND ($2::UUID IS NULL OR $2 = org_uid)
	`
	filterQuery := dbhelpers.ResolveFilterSort(dbhelpers.UserActivitiesTable, filter)
	queryStmt := fmt.Sprintf("%s %s %s", selectQuery, conditionsQuery, filterQuery)

	var queryArgs []interface{}
	queryArgs = append(queryArgs, userID, orgUID)

	// query rows
	rows, err := s.conn.Query(ctx, queryStmt, queryArgs...)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	// get total
	total, err := dbhelpers.GetGlobalTotal(ctx, s.conn, dbhelpers.UserActivitiesTable, conditionsQuery, queryArgs)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}

	result, err := s.scanRows(rows)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	return result, total, nil
}

// GetByID UserActivity
func (s *UserActivityStore) GetByID(ctx context.Context, id int64) (*dbmodels.UserActivity, *faulterr.FaultErr) {
	errMsg := "error when trying to get user by id"

	queryStmt := `SELECT * FROM user_activities WHERE id=$1`

	row := s.conn.QueryRow(ctx, queryStmt, id)
	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert UserActivity
func (s *UserActivityStore) Insert(ctx context.Context, tx pgx.Tx, arg dbmodels.UserActivity) (*dbmodels.UserActivity, *faulterr.FaultErr) {
	errMsg := "error when trying to insert user"

	queryStmt := `
	INSERT INTO
	user_activities(
		user_id,
		org_uid,
		action,
		object_id,
		object_type,
		session_token
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *
	`

	row := tx.QueryRow(ctx, queryStmt,
		&arg.UserID,
		&arg.OrgUID,
		&arg.Action,
		&arg.ObjectID,
		&arg.ObjectType,
		&arg.SessionToken,
	)

	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// Update UserActivity
func (s *UserActivityStore) Update(ctx context.Context, tx pgx.Tx, arg dbmodels.UserActivity) *faulterr.FaultErr {
	errMsg := "error when trying to update user"

	queryStmt := `
	UPDATE user_activities
	SET
		user_id=$1,
		org_uid=$2,
		action=$3,
		object_id=$4,
		object_type=$5,
		session_token=$6
	WHERE id=$7
	`

	_, err := tx.Exec(ctx, queryStmt,
		&arg.UserID,
		&arg.OrgUID,
		&arg.Action,
		&arg.ObjectID,
		&arg.ObjectType,
		&arg.SessionToken,
		&arg.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}
	return nil
}

// Delete UserActivity
func (s *UserActivityStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM user_activities WHERE id=$1`

	_, err := tx.Exec(ctx, queryStmt, id)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete user")
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *UserActivityStore) scanRows(rows pgx.Rows) ([]dbmodels.UserActivity, error) {
	result := []dbmodels.UserActivity{}
	obj := dbmodels.UserActivity{}

	for rows.Next() {
		if err := rows.Scan(
			&obj.ID,
			&obj.UserID,
			&obj.OrgUID,
			&obj.Action,
			&obj.ObjectID,
			&obj.ObjectType,
			&obj.SessionToken,
			&obj.CreatedAt,
			&obj.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, obj)
	}
	return result, nil
}

func (s *UserActivityStore) scanRow(row pgx.Row) (*dbmodels.UserActivity, error) {
	obj := &dbmodels.UserActivity{}
	if err := row.Scan(
		&obj.ID,
		&obj.UserID,
		&obj.OrgUID,
		&obj.Action,
		&obj.ObjectID,
		&obj.ObjectType,
		&obj.SessionToken,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return obj, nil
}

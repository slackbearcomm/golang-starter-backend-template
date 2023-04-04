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

type RoleStore struct {
	conn *pgxpool.Pool
}

var _ RoleStoreInterface = &RoleStore{}

type RoleStoreInterface interface {
	GetCount(ctx context.Context, orgUID uuid.UUID) (int64, *faulterr.FaultErr)
	GetManyByIDs(ctx context.Context, ids []int64) ([]*dbmodels.Role, *faulterr.FaultErr)

	List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID, deptID *int64) ([]dbmodels.Role, int, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*dbmodels.Role, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string) (*dbmodels.Role, *faulterr.FaultErr)

	Insert(ctx context.Context, tx pgx.Tx, r dbmodels.Role) (*dbmodels.Role, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, r dbmodels.Role) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewRoleStore(conn *pgxpool.Pool) *RoleStore {
	return &RoleStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetCount retrives last row from database
func (s *RoleStore) GetCount(ctx context.Context, orgUID uuid.UUID) (int64, *faulterr.FaultErr) {
	errMsg := "error getting last inserted row"

	queryStmt := `
	SELECT COUNT(*) FROM roles
	WHERE roles.org_uid = $1
	`

	var count int64
	rows, err := s.conn.Query(ctx, queryStmt, orgUID)
	if err != nil {
		return 0, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, faulterr.NewPostgresError(err, errMsg)
		}
	}
	return count, nil
}

// GetManyByIDs get all roles by ids
func (s *RoleStore) GetManyByIDs(ctx context.Context, ids []int64) ([]*dbmodels.Role, *faulterr.FaultErr) {
	errMsg := "error when trying to get many roles by ids"

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	queryStmt := "SELECT * from roles WHERE id IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := s.conn.Query(ctx, queryStmt, args...)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	output, err := s.scanRows(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	result := []*dbmodels.Role{}
	for i := 0; i < len(output); i++ {
		result = append(result, &output[i])
	}
	return result, nil
}

// List retrives all roles from database
func (s *RoleStore) List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID, deptID *int64) ([]dbmodels.Role, int, *faulterr.FaultErr) {
	errMsg := "error when trying to get roles"

	// define query
	selectQuery := `SELECT * FROM roles`
	conditionsQuery := `
	WHERE ($1::UUID IS NULL OR $1 = org_uid)
	AND ($2::INTEGER IS NULL OR $2 = department_id)
	AND ($3::BOOLEAN IS NULL OR $3 = is_final)
	AND ($4::BOOLEAN IS NULL OR $4 = is_archived)
	`
	filterQuery := dbhelpers.ResolveFilterSort(dbhelpers.RolesTable, filter)
	queryStmt := fmt.Sprintf("%s %s %s", selectQuery, conditionsQuery, filterQuery)

	var queryArgs []interface{}
	queryArgs = append(queryArgs, orgUID, deptID, filter.IsFinal, filter.IsArchived)

	// query rows
	rows, err := s.conn.Query(ctx, queryStmt, queryArgs...)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	// get total
	total, err := dbhelpers.GetGlobalTotal(ctx, s.conn, dbhelpers.RolesTable, conditionsQuery, queryArgs)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}

	result, err := s.scanRows(rows)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	return result, total, nil
}

// GetByID gets role by ID from database
func (s *RoleStore) GetByID(ctx context.Context, id int64) (*dbmodels.Role, *faulterr.FaultErr) {
	queryStmt := `
	SELECT * FROM roles
	WHERE roles.id = $1
	`
	errMsg := "error when trying to get role by id"

	row := s.conn.QueryRow(ctx, queryStmt, id)

	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// GetByCode gets role by code from database
func (s *RoleStore) GetByCode(ctx context.Context, code string) (*dbmodels.Role, *faulterr.FaultErr) {
	errMsg := "error when trying to get role by code"

	queryStmt := `
	SELECT * FROM roles
	WHERE roles.code = $1
	`

	row := s.conn.QueryRow(ctx, queryStmt, code)

	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Mutate****///////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// Insert inserts a role in database
func (s *RoleStore) Insert(ctx context.Context, tx pgx.Tx, arg dbmodels.Role) (*dbmodels.Role, *faulterr.FaultErr) {
	errMsg := "error when trying to insert role"

	queryStmt := `
	INSERT INTO
	roles(
		code,
		org_uid,
		department_id,
		name,
		permissions,
		is_management,
		status,
		is_final,
		is_archived
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING *
	`

	row := tx.QueryRow(ctx, queryStmt,
		&arg.Code,
		&arg.OrgUID,
		&arg.DepartmentID,
		&arg.Name,
		&arg.Permissions,
		&arg.IsManagement,
		&arg.Status,
		&arg.IsFinal,
		&arg.IsArchived,
	)

	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// Update updates a role in database
func (s *RoleStore) Update(ctx context.Context, tx pgx.Tx, arg dbmodels.Role) *faulterr.FaultErr {
	errMsg := "error when trying to update role"

	queryStmt := `
	UPDATE roles
	SET 
		name=$1,
		permissions=$2,
		is_management=$3,
		status=$4,
		is_final=$5,
		is_archived=$6
	WHERE id=$7
	`

	_, err := tx.Exec(ctx, queryStmt,
		&arg.Name,
		&arg.Permissions,
		&arg.IsManagement,
		&arg.Status,
		&arg.IsFinal,
		&arg.IsArchived,
		&arg.ID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}
	return nil
}

// Delete deletes a role from database
func (s *RoleStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM roles WHERE id=$1`

	_, err := tx.Exec(ctx, queryStmt)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete role")
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *RoleStore) scanRows(rows pgx.Rows) ([]dbmodels.Role, error) {
	result := []dbmodels.Role{}
	obj := dbmodels.Role{}

	for rows.Next() {
		if err := rows.Scan(
			&obj.ID,
			&obj.Code,
			&obj.OrgUID,
			&obj.DepartmentID,
			&obj.Name,
			&obj.Permissions,
			&obj.IsManagement,
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

func (s *RoleStore) scanRow(row pgx.Row) (*dbmodels.Role, error) {
	obj := dbmodels.Role{}

	if err := row.Scan(
		&obj.ID,
		&obj.Code,
		&obj.OrgUID,
		&obj.DepartmentID,
		&obj.Name,
		&obj.Permissions,
		&obj.IsManagement,
		&obj.Status,
		&obj.IsFinal,
		&obj.IsArchived,
		&obj.CreatedAt,
		&obj.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &obj, nil
}

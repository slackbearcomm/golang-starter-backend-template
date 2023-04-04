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

type DepartmentStore struct {
	conn *pgxpool.Pool
}

var _ DepartmentStoreInterface = &DepartmentStore{}

type DepartmentStoreInterface interface {
	GetCount(ctx context.Context, orgUID uuid.UUID) (int64, *faulterr.FaultErr)
	GetManyByIDs(ctx context.Context, ids []int64) ([]*dbmodels.Department, *faulterr.FaultErr)

	List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID) ([]dbmodels.Department, int, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64) (*dbmodels.Department, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string) (*dbmodels.Department, *faulterr.FaultErr)

	Insert(ctx context.Context, tx pgx.Tx, r dbmodels.Department) (*dbmodels.Department, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, r dbmodels.Department) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr
}

func NewDepartmentStore(conn *pgxpool.Pool) *DepartmentStore {
	return &DepartmentStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetCount retrives total rows count by orgUID rom database
func (s *DepartmentStore) GetCount(ctx context.Context, orgUID uuid.UUID) (int64, *faulterr.FaultErr) {
	errMsg := "error getting last inserted row"

	queryStmt := `
	SELECT COUNT(*) FROM departments
	WHERE departments.org_uid = $1
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

// GetManyByIDs get all departments by ids
func (s *DepartmentStore) GetManyByIDs(ctx context.Context, ids []int64) ([]*dbmodels.Department, *faulterr.FaultErr) {
	errMsg := "error when trying to get many departments by ids"

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i := 0; i < len(ids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = ids[i]
	}

	queryStmt := "SELECT * from departments WHERE id IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := s.conn.Query(ctx, queryStmt, args...)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	output, err := s.scanRows(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	result := []*dbmodels.Department{}
	for i := 0; i < len(output); i++ {
		result = append(result, &output[i])
	}
	return result, nil
}

// List retrives all departments from database
func (s *DepartmentStore) List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID) ([]dbmodels.Department, int, *faulterr.FaultErr) {
	errMsg := "error when trying to get departments"

	// define query
	selectQuery := `SELECT * FROM departments`
	conditionsQuery := `
	WHERE ($1::UUID IS NULL OR $1 = org_uid)
	AND ($2::BOOLEAN IS NULL OR $2 = is_final)
	AND ($3::BOOLEAN IS NULL OR $3 = is_archived)
	`
	filterQuery := dbhelpers.ResolveFilterSort(dbhelpers.DepartmentsTable, filter)
	queryStmt := fmt.Sprintf("%s %s %s", selectQuery, conditionsQuery, filterQuery)

	var queryArgs []interface{}
	queryArgs = append(queryArgs, orgUID, filter.IsFinal, filter.IsArchived)

	// query rows
	rows, err := s.conn.Query(ctx, queryStmt, queryArgs...)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	// get total
	total, err := dbhelpers.GetGlobalTotal(ctx, s.conn, dbhelpers.DepartmentsTable, conditionsQuery, queryArgs)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}

	result, err := s.scanRows(rows)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	return result, total, nil
}

// GetByID gets department by ID from database
func (s *DepartmentStore) GetByID(ctx context.Context, id int64) (*dbmodels.Department, *faulterr.FaultErr) {
	errMsg := "error when trying to get department by id"

	queryStmt := `
	SELECT * FROM departments
	WHERE departments.id = $1
	`

	row := s.conn.QueryRow(ctx, queryStmt, id)

	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// GetByCode gets department by code from database
func (s *DepartmentStore) GetByCode(ctx context.Context, code string) (*dbmodels.Department, *faulterr.FaultErr) {
	errMsg := "error when trying to get department by code"

	queryStmt := `
	SELECT * FROM departments
	WHERE departments.code = $1
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

// Insert inserts a department in database
func (s *DepartmentStore) Insert(ctx context.Context, tx pgx.Tx, arg dbmodels.Department) (*dbmodels.Department, *faulterr.FaultErr) {
	errMsg := "error when trying to insert department"

	queryStmt := `
	INSERT INTO
	departments(
		code,
		org_uid,
		name,
		status,
		is_final,
		is_archived
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *
	`

	row := tx.QueryRow(ctx, queryStmt,
		&arg.Code,
		&arg.OrgUID,
		&arg.Name,
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

// Update updates a department in database
func (s *DepartmentStore) Update(ctx context.Context, tx pgx.Tx, arg dbmodels.Department) *faulterr.FaultErr {
	errMsg := "error when trying to update department"

	queryStmt := `
	UPDATE departments
	SET 
		name=$1,
		status=$2,
		is_final=$3,
		is_archived=$4
	WHERE id=$5
	`

	_, err := tx.Exec(ctx, queryStmt,
		&arg.Name,
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

// Delete deletes a department from database
func (s *DepartmentStore) Delete(ctx context.Context, tx pgx.Tx, id int64) *faulterr.FaultErr {
	queryStmt := `DELETE FROM departments WHERE id=$1`

	_, err := tx.Exec(ctx, queryStmt)
	if err != nil {
		return faulterr.NewPostgresError(err, "error when trying to delete department")
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *DepartmentStore) scanRows(rows pgx.Rows) ([]dbmodels.Department, error) {
	result := []dbmodels.Department{}
	obj := dbmodels.Department{}

	for rows.Next() {
		if err := rows.Scan(
			&obj.ID,
			&obj.Code,
			&obj.OrgUID,
			&obj.Name,
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

func (s *DepartmentStore) scanRow(row pgx.Row) (*dbmodels.Department, error) {
	obj := dbmodels.Department{}

	if err := row.Scan(
		&obj.ID,
		&obj.Code,
		&obj.OrgUID,
		&obj.Name,
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

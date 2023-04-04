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

type OrganizationStore struct {
	conn *pgxpool.Pool
}

var _ OrganizationStoreInterface = &OrganizationStore{}

type OrganizationStoreInterface interface {
	GetCount(ctx context.Context) (int64, *faulterr.FaultErr)
	GetManyByUIDs(ctx context.Context, uids []string) ([]*dbmodels.Organization, *faulterr.FaultErr)

	List(ctx context.Context, filter models.SearchFilter, sector *string) ([]dbmodels.Organization, int, *faulterr.FaultErr)
	GetByUID(ctx context.Context, uid uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string) (*dbmodels.Organization, *faulterr.FaultErr)

	Insert(ctx context.Context, tx pgx.Tx, o dbmodels.Organization) (*dbmodels.Organization, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, o dbmodels.Organization) *faulterr.FaultErr
	Delete(ctx context.Context, tx pgx.Tx, uid uuid.UUID) *faulterr.FaultErr
}

func NewOrganizationStore(conn *pgxpool.Pool) *OrganizationStore {
	return &OrganizationStore{conn}
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Read****/////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

// GetCount retrieves last row from database
func (s *OrganizationStore) GetCount(ctx context.Context) (int64, *faulterr.FaultErr) {
	errMsg := "error getting last inserted row"

	queryStmt := `SELECT COUNT(*) FROM organizations`

	var count int64
	rows, err := s.conn.Query(ctx, queryStmt)
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

// GetManyByIDs get all organizations by ids
func (s *OrganizationStore) GetManyByUIDs(ctx context.Context, uids []string) ([]*dbmodels.Organization, *faulterr.FaultErr) {
	errMsg := "error when trying to get many organizations by uids"

	placeholders := make([]string, len(uids))
	args := make([]interface{}, len(uids))
	for i := 0; i < len(uids); i++ {
		index := strconv.Itoa(i + 1)
		placeholders[i] = "$" + index
		args[i] = uids[i]
	}

	queryStmt := "SELECT * from organizations WHERE uid IN (" + strings.Join(placeholders, ",") + ")"

	rows, err := s.conn.Query(ctx, queryStmt, args...)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	output, err := s.scanRows(rows)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}

	result := []*dbmodels.Organization{}
	for i := 0; i < len(output); i++ {
		result = append(result, &output[i])
	}
	return result, nil
}

// List retrives all organizations from database
func (s *OrganizationStore) List(ctx context.Context, filter models.SearchFilter, sector *string) ([]dbmodels.Organization, int, *faulterr.FaultErr) {
	errMsg := "error when trying to get organizations"

	// define query
	selectQuery := `SELECT * FROM organizations`
	conditionsQuery := `
	WHERE ($1::BOOLEAN IS NULL OR $1 = is_archived)
	AND ($2::VARCHAR IS NULL OR $2 = sector)
	`
	filterQuery := dbhelpers.ResolveFilterSort(dbhelpers.OrganizationsTable, filter)
	queryStmt := fmt.Sprintf("%s %s %s", selectQuery, conditionsQuery, filterQuery)

	var queryArgs []interface{}
	queryArgs = append(queryArgs, filter.IsArchived, sector)

	// query rows
	rows, err := s.conn.Query(ctx, queryStmt, queryArgs...)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	defer rows.Close()

	// get total
	total, err := dbhelpers.GetGlobalTotal(ctx, s.conn, dbhelpers.OrganizationsTable, conditionsQuery, queryArgs)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}

	result, err := s.scanRows(rows)
	if err != nil {
		return nil, 0, faulterr.NewPostgresError(err, errMsg)
	}
	return result, total, nil
}

// GetByUID gets organization by UID from database
func (s *OrganizationStore) GetByUID(ctx context.Context, uid uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr) {
	errMsg := "error when trying to get organization by uid"

	queryStmt := `
	SELECT * FROM organizations
	WHERE organizations.uid = $1
	`

	row := s.conn.QueryRow(ctx, queryStmt, uid)
	obj, err := s.scanRow(row)
	if err != nil {
		return nil, faulterr.NewPostgresError(err, errMsg)
	}
	return obj, nil
}

// GetByCode gets organization by code from database
func (s *OrganizationStore) GetByCode(ctx context.Context, code string) (*dbmodels.Organization, *faulterr.FaultErr) {
	errMsg := "error when trying to get organization by code"

	queryStmt := `
	SELECT * FROM organizations
	WHERE organizations.code = $1
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

// Insert inserts a organization in database
func (s *OrganizationStore) Insert(ctx context.Context, tx pgx.Tx, arg dbmodels.Organization) (*dbmodels.Organization, *faulterr.FaultErr) {
	errMsg := "error when trying to insert organization"

	queryStmt := `
	INSERT INTO
	organizations(
		uid,
		code,
		name,
		website,
		logo,
		sector,
		status,
		is_final,
		is_archived
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING *
	`

	row := tx.QueryRow(ctx, queryStmt,
		&arg.UID,
		&arg.Code,
		&arg.Name,
		&arg.Website,
		&arg.Logo,
		&arg.Sector,
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

// Update updates a organization in database
func (s *OrganizationStore) Update(ctx context.Context, tx pgx.Tx, arg dbmodels.Organization) *faulterr.FaultErr {
	errMsg := "error when trying to update organization"

	queryStmt := `
	UPDATE organizations
	SET
		name=$1,
		website=$2,
		logo=$3,
		sector=$4,
		status=$5,
		is_archived=$6
	WHERE uid=$7
	`

	_, err := tx.Exec(ctx, queryStmt,
		&arg.Name,
		&arg.Website,
		&arg.Logo,
		&arg.Sector,
		&arg.Status,
		&arg.IsArchived,
		&arg.UID,
	)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}
	return nil
}

// Delete deletes a organization from database
func (s *OrganizationStore) Delete(ctx context.Context, tx pgx.Tx, uid uuid.UUID) *faulterr.FaultErr {
	errMsg := "error when trying to delete organization"

	queryStmt := `DELETE FROM organizations WHERE uid=$1`

	_, err := tx.Exec(ctx, queryStmt)
	if err != nil {
		return faulterr.NewPostgresError(err, errMsg)
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////****Helpers****//////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////

func (s *OrganizationStore) scanRows(rows pgx.Rows) ([]dbmodels.Organization, error) {
	result := []dbmodels.Organization{}
	obj := dbmodels.Organization{}

	for rows.Next() {
		if err := rows.Scan(
			&obj.ID,
			&obj.UID,
			&obj.Code,
			&obj.Name,
			&obj.Website,
			&obj.Logo,
			&obj.Sector,
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

func (s *OrganizationStore) scanRow(row pgx.Row) (*dbmodels.Organization, error) {
	obj := &dbmodels.Organization{}
	if err := row.Scan(
		&obj.ID,
		&obj.UID,
		&obj.Code,
		&obj.Name,
		&obj.Website,
		&obj.Logo,
		&obj.Sector,
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

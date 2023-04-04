package orgmaster

import (
	"context"
	"gogql/app/helpers"
	"gogql/app/models/constants"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/jackc/pgx/v5"
)

type DepartmentMaster struct {
	dbstore *dbstore.DBStore
}

func NewDepartmentMaster(s *dbstore.DBStore) *DepartmentMaster {
	return &DepartmentMaster{s}
}

func (m *DepartmentMaster) VerifyOrganizationExists(ctx context.Context, id int64) (*dbmodels.Department, *faulterr.FaultErr) {
	dept, err := m.dbstore.DepartmentStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !dept.IsFinal {
		return nil, faulterr.NewFrobiddenError("department is not yet final")
	}
	if dept.IsArchived {
		return nil, faulterr.NewFrobiddenError("department is archived")
	}
	return dept, nil
}

func (m *DepartmentMaster) CreateOne(ctx context.Context, tx pgx.Tx, req dbmodels.DepartmentRequest, org dbmodels.Organization) (*dbmodels.Department, *faulterr.FaultErr) {
	requests := []dbmodels.DepartmentRequest{}
	requests = append(requests, req)

	result, err := m.BulkCreate(ctx, tx, requests, org)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func (m *DepartmentMaster) BulkCreate(ctx context.Context, tx pgx.Tx, requests []dbmodels.DepartmentRequest, org dbmodels.Organization) ([]*dbmodels.Department, *faulterr.FaultErr) {
	result := []*dbmodels.Department{}

	// Get count
	count, err := m.dbstore.DepartmentStore.GetCount(ctx, org.UID)
	if err != nil {
		return nil, err
	}

	for i := range requests {
		// validate request
		if err := m.validate(requests[i]); err != nil {
			return nil, err
		}

		// construct arguments
		arg := m.construct(requests[i])
		arg.Code = helpers.GenerateCode(constants.DepartmentObject, org.Code, count)
		if arg.Status == "" {
			arg.Status = constants.StatusCreated
		}

		// insert into db
		obj, err := m.dbstore.DepartmentStore.Insert(ctx, tx, *arg)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
		count++
	}

	return result, nil
}

func (m *DepartmentMaster) construct(req dbmodels.DepartmentRequest) *dbmodels.Department {
	return &dbmodels.Department{
		OrgUID:     req.OrgUID,
		Name:       req.Name,
		IsFinal:    req.IsFinal,
		IsArchived: false,
	}
}

func (m *DepartmentMaster) validate(req dbmodels.DepartmentRequest) *faulterr.FaultErr {
	if req.OrgUID.String() == "" {
		return faulterr.NewBadRequestError("organization uid is required")
	}
	if req.Name == "" {
		return faulterr.NewBadRequestError("department name is required")
	}
	return nil
}

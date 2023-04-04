package orgservice

import (
	"context"
	"gogql/app/master"
	"gogql/app/models"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type DepartmentService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ DepartmentServiceInterface = &DepartmentService{}

type DepartmentServiceInterface interface {
	List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID) ([]dbmodels.Department, int, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr)

	Create(ctx context.Context, tx pgx.Tx, req dbmodels.DepartmentRequest) (*dbmodels.Department, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, id int64, req dbmodels.DepartmentRequest, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr)
	Delete(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) *faulterr.FaultErr
}

func NewDepartmentService(s *dbstore.DBStore, m *master.Master) *DepartmentService {
	return &DepartmentService{s, m}
}

// List gets all departments for super admin and associated organization departments for members
func (s *DepartmentService) List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID) ([]dbmodels.Department, int, *faulterr.FaultErr) {
	return s.dbstore.DepartmentStore.List(ctx, filter, orgUID)
}

// GetByID gets a department by department id
func (s *DepartmentService) GetByID(ctx context.Context, id int64, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr) {
	obj, err := s.dbstore.DepartmentStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if orgUID != nil && *orgUID != obj.OrgUID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

// GetByID gets a department by department code
func (s *DepartmentService) GetByCode(ctx context.Context, code string, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr) {
	obj, err := s.dbstore.DepartmentStore.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if orgUID != nil && *orgUID != obj.OrgUID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

// Create saves a department object in db
func (s *DepartmentService) Create(ctx context.Context, tx pgx.Tx, req dbmodels.DepartmentRequest) (*dbmodels.Department, *faulterr.FaultErr) {
	// verify if org exists
	org, err := s.master.OrganizationMaster.VerifyOrganizationExists(ctx, req.OrgUID)
	if err != nil {
		return nil, err
	}
	return s.master.DepartmentMaster.CreateOne(ctx, tx, req, *org)
}

// CreateMany saves a department object in db
func (s *DepartmentService) CreateMany(ctx context.Context, tx pgx.Tx, req []dbmodels.DepartmentRequest, org dbmodels.Organization) ([]*dbmodels.Department, *faulterr.FaultErr) {
	return s.master.DepartmentMaster.BulkCreate(ctx, tx, req, org)
}

func (s *DepartmentService) Update(ctx context.Context, tx pgx.Tx, id int64, req dbmodels.DepartmentRequest, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}

	// update fields
	if req.Name != "" {
		obj.Name = req.Name
	}

	// update department
	if err := s.dbstore.DepartmentStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *DepartmentService) Finalize(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	if obj.IsFinal {
		return nil, faulterr.NewBadRequestError("department is already final")
	}

	obj.IsFinal = true
	if err := s.dbstore.DepartmentStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *DepartmentService) Archive(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	if obj.IsArchived {
		return nil, faulterr.NewBadRequestError("department is already archived")
	}

	obj.IsArchived = true
	if err := s.dbstore.DepartmentStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *DepartmentService) Unarchive(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) (*dbmodels.Department, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	if !obj.IsArchived {
		return nil, faulterr.NewBadRequestError("department is already unarchived")
	}

	obj.IsArchived = false
	if err := s.dbstore.DepartmentStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *DepartmentService) Delete(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) *faulterr.FaultErr {
	_, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return err
	}
	return s.dbstore.DepartmentStore.Delete(ctx, tx, id)
}

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

type RoleService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ RoleServiceInterface = &RoleService{}

type RoleServiceInterface interface {
	List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID, deptID *int64) ([]dbmodels.Role, int, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr)
	Create(ctx context.Context, tx pgx.Tx, req dbmodels.RoleRequest) (*dbmodels.Role, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, id int64, req dbmodels.RoleRequest, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr)
	Delete(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) *faulterr.FaultErr
}

func NewRoleService(s *dbstore.DBStore, m *master.Master) *RoleService {
	return &RoleService{s, m}
}

// List gets all roles for super admin and associated organization roles for members
func (s *RoleService) List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID, deptID *int64) ([]dbmodels.Role, int, *faulterr.FaultErr) {
	return s.dbstore.RoleStore.List(ctx, filter, orgUID, deptID)
}

// GetByID gets a role by role id
func (s *RoleService) GetByID(ctx context.Context, id int64, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr) {
	obj, err := s.dbstore.RoleStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if orgUID != nil && *orgUID != obj.OrgUID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

// GetByID gets a role by role code
func (s *RoleService) GetByCode(ctx context.Context, code string, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr) {
	obj, err := s.dbstore.RoleStore.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if orgUID != nil && *orgUID != obj.OrgUID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

// Create saves a role object in db
func (s *RoleService) Create(ctx context.Context, tx pgx.Tx, req dbmodels.RoleRequest) (*dbmodels.Role, *faulterr.FaultErr) {
	// verify organization
	org, err := s.master.OrganizationMaster.VerifyOrganizationExists(ctx, req.OrgUID)
	if err != nil {
		return nil, err
	}
	// verify department
	_, err = s.master.DepartmentMaster.VerifyOrganizationExists(ctx, req.DepartmentID)
	if err != nil {
		return nil, err
	}

	// create role
	role, err := s.master.RoleMaster.CreateOne(ctx, tx, req, *org)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) Update(ctx context.Context, tx pgx.Tx, id int64, req dbmodels.RoleRequest, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}

	// update fields
	if req.Name != "" {
		obj.Name = req.Name
	}
	obj.IsManagement = req.IsManagement
	obj.Permissions = req.Permissions
	obj.IsFinal = req.IsFinal

	// update role
	if err := s.dbstore.RoleStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *RoleService) Finalize(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	if obj.IsFinal {
		return nil, faulterr.NewBadRequestError("role is already final")
	}

	obj.IsFinal = true
	if err := s.dbstore.RoleStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *RoleService) Archive(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	if obj.IsArchived {
		return nil, faulterr.NewBadRequestError("role is already archived")
	}

	obj.IsArchived = true
	if err := s.dbstore.RoleStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *RoleService) Unarchive(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) (*dbmodels.Role, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	if !obj.IsArchived {
		return nil, faulterr.NewBadRequestError("role is already unarchived")
	}

	obj.IsArchived = false
	if err := s.dbstore.RoleStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *RoleService) Delete(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) *faulterr.FaultErr {
	_, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return err
	}
	return s.dbstore.RoleStore.Delete(ctx, tx, id)
}

// Unique permissions
func (s *RoleService) UniquePermissions(list []string) []string {
	permissions := []string{}

	for _, perm := range list {
		if !s.ContainsPermission(permissions[:], perm) {
			permissions = append(permissions, perm)
		}
	}
	return permissions
}

func (s *RoleService) ContainsPermission(list []string, perm string) bool {
	for _, e := range list {
		if e == perm {
			return true
		}
	}
	return false
}

package orgservice

import (
	"context"
	"gogql/app/helpers"
	"gogql/app/master"
	"gogql/app/models"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/volatiletech/null"
)

type OrganizationService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ OrganizationServiceInterface = &OrganizationService{}

type OrganizationServiceInterface interface {
	List(ctx context.Context, filter models.SearchFilter, sector *string) ([]dbmodels.Organization, int, *faulterr.FaultErr)
	GetByUID(ctx context.Context, uid uuid.UUID, orgUID *uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr)
	GetByCode(ctx context.Context, code string, orgUID *uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, uid uuid.UUID, req dbmodels.OrganizationRequest, orgUID *uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr)
	Archive(ctx context.Context, tx pgx.Tx, uid uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr)
	Unarchive(ctx context.Context, tx pgx.Tx, uid uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr)
	Delete(ctx context.Context, tx pgx.Tx, uid uuid.UUID) *faulterr.FaultErr
}

func NewOrganizationService(s *dbstore.DBStore, m *master.Master) *OrganizationService {
	return &OrganizationService{s, m}
}

// List gets all skus
func (s *OrganizationService) List(ctx context.Context, filter models.SearchFilter, sector *string) ([]dbmodels.Organization, int, *faulterr.FaultErr) {
	return s.dbstore.OrganizationStore.List(ctx, filter, sector)
}

func (s *OrganizationService) GetByUID(ctx context.Context, uid uuid.UUID, orgUID *uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr) {
	obj, err := s.dbstore.OrganizationStore.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if orgUID != nil && *orgUID != obj.UID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

func (s *OrganizationService) GetByCode(ctx context.Context, code string, orgUID *uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr) {
	obj, err := s.dbstore.OrganizationStore.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if orgUID != nil && *orgUID != obj.UID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

// Create inserts organization, role, and a user for the organization
func (s *OrganizationService) Create(ctx context.Context, tx pgx.Tx, req dbmodels.OrganizationRegisterRequest) (*dbmodels.Organization, *dbmodels.User, *faulterr.FaultErr) {
	orgReq := dbmodels.OrganizationRequest{
		Name:    req.OrgName,
		Website: req.Website,
		Logo:    req.Logo,
		Sector:  req.Sector,
	}

	org, err := s.master.OrganizationMaster.CreateOne(ctx, tx, orgReq)
	if err != nil {
		return nil, nil, err
	}

	// create management department
	deptReq := dbmodels.DepartmentRequest{
		OrgUID:  org.UID,
		Name:    models.RoleManagementStr,
		IsFinal: true,
	}

	dept, err := s.master.DepartmentMaster.CreateOne(ctx, tx, deptReq, *org)
	if err != nil {
		return nil, nil, err
	}

	// Insert management role
	roleReq := dbmodels.RoleRequest{
		OrgUID:       org.UID,
		DepartmentID: dept.ID,
		Name:         models.RoleManagementStr,
		Permissions:  nil,
		IsManagement: true,
	}

	role, err := s.master.RoleMaster.CreateOne(ctx, tx, roleReq, *org)
	if err != nil {
		return nil, nil, err
	}

	// Insert user with management role
	userReq := dbmodels.UserRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		OrgUID:    helpers.NullUUIDFromUUID(org.UID),
		RoleID:    null.Int64From(role.ID),
	}

	user, err := s.master.UserMaster.CreateOne(ctx, tx, userReq)
	if err != nil {
		return nil, nil, err
	}

	// push to nexport
	// _, err = s.rest.OrganizationStore.Create(org)
	// if err != nil {
	// 	return nil, nil, err
	// }

	return org, user, nil
}

func (s *OrganizationService) Update(ctx context.Context, tx pgx.Tx, uid uuid.UUID, req dbmodels.OrganizationRequest, orgUID *uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr) {
	obj, err := s.GetByUID(ctx, uid, orgUID)
	if err != nil {
		return nil, err
	}

	// Compare changes
	if req.Name != "" {
		obj.Name = req.Name
	}
	obj.Website = req.Website

	if err := s.dbstore.OrganizationStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *OrganizationService) Archive(ctx context.Context, tx pgx.Tx, uid uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr) {
	obj, err := s.GetByUID(ctx, uid, nil)
	if err != nil {
		return nil, err
	}
	if obj.IsArchived {
		return nil, faulterr.NewBadRequestError("organization is already archived")
	}

	obj.IsArchived = true
	if err := s.dbstore.OrganizationStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *OrganizationService) Unarchive(ctx context.Context, tx pgx.Tx, uid uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr) {
	obj, err := s.GetByUID(ctx, uid, nil)
	if err != nil {
		return nil, err
	}
	if !obj.IsArchived {
		return nil, faulterr.NewBadRequestError("organization is already unarchived")
	}

	obj.IsArchived = false
	if err := s.dbstore.OrganizationStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func (s *OrganizationService) Delete(ctx context.Context, tx pgx.Tx, uid uuid.UUID) *faulterr.FaultErr {
	_, err := s.dbstore.OrganizationStore.GetByUID(ctx, uid)
	if err != nil {
		return err
	}
	return s.dbstore.OrganizationStore.Delete(ctx, tx, uid)
}

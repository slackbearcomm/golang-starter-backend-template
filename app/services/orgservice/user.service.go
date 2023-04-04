package orgservice

import (
	"context"
	"gogql/app/master"
	"gogql/app/models"
	"gogql/app/models/constants"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ UserServiceInterface = &UserService{}

type UserServiceInterface interface {
	Me(ctx context.Context, userID int64) (*dbmodels.User, *faulterr.FaultErr)
	List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID, roleID *int64) ([]dbmodels.User, int, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, orgUID *uuid.UUID) (*dbmodels.User, *faulterr.FaultErr)
	Update(ctx context.Context, tx pgx.Tx, id int64, request dbmodels.UserRequest, orgUID *uuid.UUID) (*dbmodels.User, *faulterr.FaultErr)
	Delete(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) *faulterr.FaultErr
}

func NewUserService(s *dbstore.DBStore, m *master.Master) *UserService {
	return &UserService{s, m}
}

// ListCustomers gets aa user by id
func (s *UserService) Me(ctx context.Context, userID int64) (*dbmodels.User, *faulterr.FaultErr) {
	return s.dbstore.UserStore.GetByID(ctx, userID)
}

// List gets all admin, members, and consumers
func (s *UserService) List(ctx context.Context, filter models.SearchFilter, orgUID *uuid.UUID, roleID *int64) ([]dbmodels.User, int, *faulterr.FaultErr) {
	return s.dbstore.UserStore.List(ctx, filter, orgUID, roleID)
}

// GetByID gets a user by its ID
func (s *UserService) GetByID(ctx context.Context, id int64, orgUID *uuid.UUID) (*dbmodels.User, *faulterr.FaultErr) {
	obj, err := s.dbstore.UserStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if orgUID != nil && obj.OrgUID.Valid && *orgUID != obj.OrgUID.UUID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

// GetByEmail gets a user by user profile
func (s *UserService) GetByEmail(ctx context.Context, email string, orgUID *uuid.UUID) (*dbmodels.User, *faulterr.FaultErr) {
	obj, err := s.dbstore.UserStore.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if orgUID != nil && obj.OrgUID.Valid && *orgUID != obj.OrgUID.UUID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

// GetByPhone gets a user by user profile
func (s *UserService) GetByPhone(ctx context.Context, phone string, orgUID *uuid.UUID) (*dbmodels.User, *faulterr.FaultErr) {
	obj, err := s.dbstore.UserStore.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if orgUID != nil && obj.OrgUID.Valid && *orgUID != obj.OrgUID.UUID {
		return nil, faulterr.NewNotFoundError("object not found")
	}

	return obj, nil
}

// Create saves a user object in db
func (s *UserService) Create(ctx context.Context, tx pgx.Tx, request dbmodels.UserRequest) (*dbmodels.User, *faulterr.FaultErr) {
	return s.master.UserMaster.CreateOne(ctx, tx, request)
}

func (s *UserService) Update(ctx context.Context, tx pgx.Tx, id int64, request dbmodels.UserRequest, orgUID *uuid.UUID) (*dbmodels.User, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	return s.master.UserMaster.Update(ctx, tx, *obj, request)
}

func (s *UserService) Delete(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) *faulterr.FaultErr {
	_, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return err
	}
	return s.dbstore.UserStore.Delete(ctx, tx, id)
}

// Archive updates a user object in db
func (s *UserService) Archive(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) (*dbmodels.User, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	if obj.IsArchived {
		return nil, faulterr.NewBadRequestError("user is already archived")
	}

	obj.IsArchived = true
	obj.Status = constants.StatusArchived
	if err := s.dbstore.UserStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// Unarchive updates a user object in db
func (s *UserService) Unarchive(ctx context.Context, tx pgx.Tx, id int64, orgUID *uuid.UUID) (*dbmodels.User, *faulterr.FaultErr) {
	obj, err := s.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err
	}
	if !obj.IsArchived {
		return nil, faulterr.NewBadRequestError("user is already unarchived")
	}

	obj.IsArchived = false
	obj.Status = constants.StatusActive
	if err := s.dbstore.UserStore.Update(ctx, tx, *obj); err != nil {
		return nil, err
	}

	return obj, nil
}

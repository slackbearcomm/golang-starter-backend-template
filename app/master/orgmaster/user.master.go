package orgmaster

import (
	"context"
	"gogql/app/models/constants"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/jackc/pgx/v5"
)

type UserMaster struct {
	dbstore *dbstore.DBStore
}

func NewUserMaster(dbstore *dbstore.DBStore) *UserMaster {
	return &UserMaster{dbstore}
}

// CreateUser creates and saves user in the db
func (m *UserMaster) CreateOne(ctx context.Context, tx pgx.Tx, r dbmodels.UserRequest) (*dbmodels.User, *faulterr.FaultErr) {
	if err := m.validate(r); err != nil {
		return nil, err
	}

	obj := dbmodels.User{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Phone:     r.Phone,
		OrgUID:    r.OrgUID,
		RoleID:    r.RoleID,
		IsAdmin:   r.IsAdmin,
		IsFinal:   true,
		Status:    constants.StatusActive,
	}

	err := m.verifyUniqueFields(ctx, obj)
	if err != nil {
		return nil, err
	}

	return m.dbstore.UserStore.Insert(ctx, tx, obj)
}

func (s *UserMaster) Update(ctx context.Context, tx pgx.Tx, obj dbmodels.User, req dbmodels.UserRequest) (*dbmodels.User, *faulterr.FaultErr) {
	// Verify request fields
	if req.FirstName != "" {
		obj.FirstName = req.FirstName
	}
	if req.LastName != "" {
		obj.LastName = req.LastName
	}
	if req.Email != "" && req.Email != obj.Email {
		// verify unique email
		_, err := s.dbstore.UserStore.GetByEmail(ctx, req.Email)
		if err == nil {
			return nil, faulterr.NewBadRequestError("email already registered")
		}
		obj.Email = req.Email
	}
	if req.Phone != "" && req.Phone != obj.Phone {
		// verify unique phone
		_, err := s.dbstore.UserStore.GetByPhone(ctx, req.Phone)
		if err == nil {
			return nil, faulterr.NewBadRequestError("phone already registered")
		}
		obj.Phone = req.Phone
	}

	return &obj, nil
}

// Validators

// verifyUniqueFields verifies the uniqueness of user
func (m *UserMaster) verifyUniqueFields(ctx context.Context, u dbmodels.User) *faulterr.FaultErr {
	// Verify unique email
	_, err := m.dbstore.UserStore.GetByEmail(ctx, u.Email)
	if err == nil {
		return faulterr.NewBadRequestError("email already registered")
	}

	// Verify unique phone
	_, err = m.dbstore.UserStore.GetByPhone(ctx, u.Phone)
	if err == nil {
		return faulterr.NewBadRequestError("phone already registered")
	}

	return nil
}

// validate verifies the uniqueness of user
func (m *UserMaster) validate(r dbmodels.UserRequest) *faulterr.FaultErr {
	if r.FirstName == "" {
		return faulterr.NewBadRequestError("First Name is required")
	}
	if r.LastName == "" {
		return faulterr.NewBadRequestError("Last Name is required")
	}
	if r.Email == "" {
		return faulterr.NewBadRequestError("Email is required")
	}

	return nil
}

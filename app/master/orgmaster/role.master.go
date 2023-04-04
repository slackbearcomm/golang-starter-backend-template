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

type RoleMaster struct {
	dbstore *dbstore.DBStore
}

func NewRoleMaster(s *dbstore.DBStore) *RoleMaster {
	return &RoleMaster{s}
}

func (m *RoleMaster) CreateOne(ctx context.Context, tx pgx.Tx, req dbmodels.RoleRequest, org dbmodels.Organization) (*dbmodels.Role, *faulterr.FaultErr) {
	requests := []dbmodels.RoleRequest{}
	requests = append(requests, req)

	result, err := m.BulkCreate(ctx, tx, requests, org)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func (m *RoleMaster) BulkCreate(ctx context.Context, tx pgx.Tx, requests []dbmodels.RoleRequest, org dbmodels.Organization) ([]*dbmodels.Role, *faulterr.FaultErr) {
	result := []*dbmodels.Role{}

	// Get count
	count, err := m.dbstore.RoleStore.GetCount(ctx, org.UID)
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
		arg.Code = helpers.GenerateCode(constants.RoleObject, org.Code, count)
		if arg.Status == "" {
			arg.Status = constants.StatusCreated
		}

		// insert into db
		obj, err := m.dbstore.RoleStore.Insert(ctx, tx, *arg)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
		count++
	}

	return result, nil
}

func (m *RoleMaster) construct(r dbmodels.RoleRequest) *dbmodels.Role {
	return &dbmodels.Role{
		OrgUID:       r.OrgUID,
		DepartmentID: r.DepartmentID,
		Name:         r.Name,
		Permissions:  r.Permissions,
		IsManagement: r.IsManagement,
		IsFinal:      r.IsFinal,
		IsArchived:   false,
	}
}

func (m *RoleMaster) GrantPermission(ctx context.Context, roleID int64, permission string) *faulterr.FaultErr {
	role, err := m.dbstore.RoleStore.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	if !role.IsManagement {
		for _, perm := range role.Permissions {
			if perm == permission {
				return nil
			}
		}
		return faulterr.NewUnauthorizedError("User Not Authorized")
	}

	return nil
}

func (m *RoleMaster) validate(r dbmodels.RoleRequest) *faulterr.FaultErr {
	if r.Name == "" {
		return faulterr.NewBadRequestError("Role Name is required")
	}
	// TODO: Review permissions condition
	// if !r.IsManagement && len(r.Permissions) == 0 {
	// 	return faulterr.NewBadRequestError("Permissions cannot be empty")
	// }
	if r.OrgUID.String() == "" {
		return faulterr.NewBadRequestError("Organization UID is required")
	}
	return nil
}

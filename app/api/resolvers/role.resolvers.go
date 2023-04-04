package resolvers

import (
	"context"
	"fmt"
	"gogql/app/api/dataloaders"
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/middlewares"
	"gogql/app/models"
	"gogql/app/models/constants"
	"gogql/app/models/dbmodels"
	"gogql/utils/faulterr"

	"github.com/volatiletech/null"
)

type roleResolver struct{ *Resolver }

// Role returns graph.RoleResolver implementation.
func (r *Resolver) Role() graph.RoleResolver { return &roleResolver{r} }

// Organization is the resolver for the organization field.
func (r *roleResolver) Organization(ctx context.Context, obj *dbmodels.Role) (*dbmodels.Organization, error) {
	return dataloaders.OrganizationLoaderFromContext(ctx, obj.OrgUID.String())
}

// Department is the resolver for the department field.
func (r *roleResolver) Department(ctx context.Context, obj *dbmodels.Role) (*dbmodels.Department, error) {
	return dataloaders.DepartmentLoaderFromContext(ctx, obj.DepartmentID)
}

///////////////
//   Query   //
///////////////

// Roles is the resolver for the roles field.
func (r *queryResolver) Roles(ctx context.Context, search graph.SearchFilter, deptID *int64) (*graph.RolesResult, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	if search.OrgUID != nil {
		orgUID = search.OrgUID
	}

	auther, err := r.GetAutherWithPermission(ctx, models.ReadRole)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	filter := r.SearchFilter(search)
	output, total, err := r.services.RoleService.List(ctx, filter, orgUID, deptID)
	if err != nil {
		return nil, err.Error
	}
	return &graph.RolesResult{Roles: output, Total: total}, nil
}

// Role is the resolver for the role field.
func (r *queryResolver) Role(ctx context.Context, id *int64, code *string) (*dbmodels.Role, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.ReadRole)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	if id != nil {
		obj, err := r.services.RoleService.GetByID(ctx, *id, orgUID)
		if err != nil {
			return nil, err.Error
		}
		return obj, nil
	}

	if code != nil {
		obj, err := r.services.RoleService.GetByCode(ctx, *code, orgUID)
		if err != nil {
			return nil, err.Error
		}
		return obj, nil
	}

	return nil, faulterr.NewFrobiddenError(noQueryParamsErr).Error
}

///////////////
// Mutations //
///////////////

// RoleCreate is the resolver for the roleCreate field.
func (r *mutationResolver) RoleCreate(ctx context.Context, input graph.UpdateRole) (*dbmodels.Role, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.CreateRole)
	if err != nil {
		return nil, err.Error
	}

	req, reqErr := r.generateRoleRequest(input)
	if reqErr != nil {
		return nil, reqErr
	}

	if auther.IsAdmin {
		if input.OrgUID != nil && input.OrgUID.Valid {
			req.OrgUID = input.OrgUID.UUID
		} else if orgUID != nil {
			req.OrgUID = *orgUID
		} else {
			return nil, faulterr.NewBadRequestError("org uid is required").Error
		}
	} else {
		req.OrgUID = auther.OrgUID.UUID
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.RoleService.Create(ctx, tx, *req)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.RoleObject, constants.CreateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.RoleObject)),
		SessionToken: auther.SessionToken,
	}
	_, err = r.services.UserActivityService.Create(ctx, tx, actReq)
	if err != nil {
		return nil, err.Error
	}

	// commit db transaction
	if err := r.services.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err.Error
	}

	return obj, nil
}

// RoleUpdate is the resolver for the roleUpdate field.
func (r *mutationResolver) RoleUpdate(ctx context.Context, id int64, input graph.UpdateRole) (*dbmodels.Role, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateRole)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	req, reqErr := r.generateRoleRequest(input)
	if reqErr != nil {
		return nil, reqErr
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.RoleService.Update(ctx, tx, id, *req, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.RoleObject, constants.UpdateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.RoleObject)),
		SessionToken: auther.SessionToken,
	}
	_, err = r.services.UserActivityService.Create(ctx, tx, actReq)
	if err != nil {
		return nil, err.Error
	}

	// commit db transaction
	if err := r.services.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err.Error
	}

	return obj, nil
}

// RoleFinalize is the resolver for the roleFinalize field.
func (r *mutationResolver) RoleFinalize(ctx context.Context, id int64) (*dbmodels.Role, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateRole)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.RoleService.Finalize(ctx, tx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.RoleObject, constants.FinalizeAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.RoleObject)),
		SessionToken: auther.SessionToken,
	}
	_, err = r.services.UserActivityService.Create(ctx, tx, actReq)
	if err != nil {
		return nil, err.Error
	}

	// commit db transaction
	if err := r.services.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err.Error
	}
	return obj, nil
}

// RoleArchive is the resolver for the roleArchive field.
func (r *mutationResolver) RoleArchive(ctx context.Context, id int64) (*dbmodels.Role, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateRole)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.RoleService.Archive(ctx, tx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.RoleObject, constants.ArchiveAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.RoleObject)),
		SessionToken: auther.SessionToken,
	}
	_, err = r.services.UserActivityService.Create(ctx, tx, actReq)
	if err != nil {
		return nil, err.Error
	}

	// commit db transaction
	if err := r.services.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err.Error
	}
	return obj, nil
}

// RoleUnarchive is the resolver for the roleUnarchive field.
func (r *mutationResolver) RoleUnarchive(ctx context.Context, id int64) (*dbmodels.Role, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateRole)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.RoleService.Unarchive(ctx, tx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.RoleObject, constants.UnarchiveAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.RoleObject)),
		SessionToken: auther.SessionToken,
	}
	_, err = r.services.UserActivityService.Create(ctx, tx, actReq)
	if err != nil {
		return nil, err.Error
	}

	// commit db transaction
	if err := r.services.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err.Error
	}
	return obj, nil
}

func (r *mutationResolver) generateRoleRequest(input graph.UpdateRole) (*dbmodels.RoleRequest, error) {
	req := &dbmodels.RoleRequest{}
	rolePermissions := []string{}

	if input.DepartmentID == nil || !input.DepartmentID.Valid || input.DepartmentID.Int64 == 0 {
		return nil, fmt.Errorf("department id is required")
	} else {
		req.DepartmentID = input.DepartmentID.Int64
	}

	if input.Name == nil || !input.Name.Valid || input.Name.String == "" {
		return nil, fmt.Errorf("role name is required")
	} else {
		req.Name = input.Name.String
	}

	if input.IsManagement == nil || !input.IsManagement.Valid {
		req.IsManagement = false
	} else {
		req.IsManagement = input.IsManagement.Bool
	}

	if !req.IsManagement && len(input.Permissions) > 0 {
		permissions := models.ListPermissions()
		reqPermissions := r.services.RoleService.UniquePermissions(input.Permissions)
		for _, perm := range reqPermissions {
			ok := r.services.RoleService.ContainsPermission(permissions, perm)
			if ok {
				rolePermissions = append(rolePermissions, perm)
			}
		}
		req.Permissions = rolePermissions
	}

	return req, nil
}

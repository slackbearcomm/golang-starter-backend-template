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

// User returns graph.UserResolver implementation.
func (r *Resolver) User() graph.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }

// Organization is the resolver for the organization field.
func (r *userResolver) Organization(ctx context.Context, obj *dbmodels.User) (*dbmodels.Organization, error) {
	if obj.OrgUID.Valid {
		return dataloaders.OrganizationLoaderFromContext(ctx, obj.OrgUID.UUID.String())
	}
	return nil, nil
}

// Role is the resolver for the role field.
func (r *userResolver) Role(ctx context.Context, obj *dbmodels.User) (*dbmodels.Role, error) {
	if obj.RoleID.Valid {
		return dataloaders.RoleLoaderFromContext(ctx, obj.RoleID.Int64)
	}
	return nil, nil
}

///////////////
//   Query   //
///////////////

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*dbmodels.User, error) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}

	obj, err := r.services.UserService.Me(ctx, auther.ID)
	if err != nil {
		return nil, err.Error
	}

	return obj, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, search graph.SearchFilter, roleID *int64) (*graph.UserResult, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	if search.OrgUID != nil {
		orgUID = search.OrgUID
	}

	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadUser); err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	filter := r.SearchFilter(search)
	output, total, err := r.services.UserService.List(ctx, filter, orgUID, roleID)
	if err != nil {
		return nil, err.Error
	}
	return &graph.UserResult{Users: output, Total: total}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id *int64, email, phone *string) (*dbmodels.User, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.ReadUser); err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	if id != nil {
		result, err := r.services.UserService.GetByID(ctx, *id, orgUID)
		if err != nil {
			return nil, err.Error
		}
		return result, nil
	}

	if email != nil {
		result, err := r.services.UserService.GetByEmail(ctx, *email, orgUID)
		if err != nil {
			return nil, err.Error
		}
		return result, nil
	}

	if phone != nil {
		result, err := r.services.UserService.GetByPhone(ctx, *phone, orgUID)
		if err != nil {
			return nil, err.Error
		}
		return result, nil
	}

	return nil, faulterr.NewFrobiddenError(noQueryParamsErr).Error
}

///////////////
// Mutations //
///////////////

// SuperAdminCreate is the resolver for the superAdminCreate field.
func (r *mutationResolver) SuperAdminCreate(ctx context.Context, input graph.UpdateUser) (*dbmodels.User, error) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		return nil, faulterr.NewUnauthorizedError(unauthorizedErr).Error
	}

	req, err := r.generateUserRequest(input)
	if err != nil {
		return nil, err.Error
	}
	req.IsAdmin = true

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.UserService.Create(ctx, tx, *req)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.UserObject, constants.CreateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.UserObject)),
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

// UserCreate is the resolver for the userCreate field.
func (r *mutationResolver) UserCreate(ctx context.Context, input graph.UpdateUser) (*dbmodels.User, error) {
	auther, err := r.GetAutherWithPermission(ctx, models.CreateUser)
	if err != nil {
		return nil, err.Error
	}

	req, err := r.generateUserRequest(input)
	if err != nil {
		return nil, err.Error
	}
	if input.OrgUID != nil && input.OrgUID.Valid {
		req.OrgUID = *input.OrgUID
	} else {
		return nil, faulterr.NewFrobiddenError("organization uid is required").Error
	}
	if input.RoleID != nil && input.RoleID.Int64 > 0 {
		req.RoleID = *input.RoleID
	} else {
		return nil, faulterr.NewFrobiddenError("role id is required").Error
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.UserService.Create(ctx, tx, *req)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.UserObject, constants.CreateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.UserObject)),
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

// ChangeDetails is the resolver for the changeDetails field.
func (r *mutationResolver) ChangeDetails(ctx context.Context, id int64, input graph.UpdateUser) (*dbmodels.User, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	req, err := r.generateUserRequest(input)
	if err != nil {
		return nil, err.Error
	}

	if input.RoleID != nil && input.RoleID.Int64 > 0 {
		req.RoleID = *input.RoleID
	}
	if auther.IsAdmin {
		if input.OrgUID != nil && input.OrgUID.Valid {
			req.OrgUID = *input.OrgUID
		} else {
			return nil, faulterr.NewBadRequestError("org uid is required").Error
		}
	} else {
		req.OrgUID = auther.OrgUID
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.UserService.Update(ctx, tx, id, *req, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.UserObject, constants.UpdateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.UserObject)),
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

// UserUpdate is the resolver for the userUpdate field.
func (r *mutationResolver) UserUpdate(ctx context.Context, id int64, input graph.UpdateUser) (*dbmodels.User, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	req, err := r.generateUserRequest(input)
	if err != nil {
		return nil, err.Error
	}

	if input.RoleID != nil && input.RoleID.Int64 > 0 {
		req.RoleID = *input.RoleID
	}
	if auther.IsAdmin {
		if input.OrgUID == nil || !input.OrgUID.Valid {
			return nil, faulterr.NewFrobiddenError("organization id is required").Error
		} else {
			req.OrgUID = *input.OrgUID
		}
	}
	if !auther.IsAdmin {
		req.OrgUID = auther.OrgUID
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.UserService.Update(ctx, tx, id, *req, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.UserObject, constants.UpdateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.UserObject)),
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

// ResendEmailVerification is the resolver for the resendEmailVerification field.
func (r *mutationResolver) ResendEmailVerification(ctx context.Context, email string) (bool, error) {
	return false, faulterr.NewBadRequestError("not implemented").Error
}

// UserArchive is the resolver for the userArchive field.
func (r *mutationResolver) UserArchive(ctx context.Context, id int64) (*dbmodels.User, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateUser)
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

	obj, err := r.services.UserService.Archive(ctx, tx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.UserObject, constants.ArchiveAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.UserObject)),
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

// UserUnarchive is the resolver for the userUnarchive field.
func (r *mutationResolver) UserUnarchive(ctx context.Context, id int64) (*dbmodels.User, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateUser)
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

	obj, err := r.services.UserService.Unarchive(ctx, tx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.UserObject, constants.UnarchiveAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.UserObject)),
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

func (r *mutationResolver) generateUserRequest(input graph.UpdateUser) (*dbmodels.UserRequest, *faulterr.FaultErr) {
	req := &dbmodels.UserRequest{}

	if input.FirstName != nil && input.FirstName.String != "" {
		req.FirstName = input.FirstName.String
	} else {
		return nil, faulterr.NewFrobiddenError("first name is required")
	}
	if input.LastName != nil && input.LastName.String != "" {
		req.LastName = input.LastName.String
	} else {
		return nil, faulterr.NewFrobiddenError("last name is required")
	}
	if input.Email != nil && input.Email.String != "" {
		req.Email = input.Email.String
	} else {
		return nil, faulterr.NewFrobiddenError("email is required")
	}
	if input.Phone != nil && input.Phone.String != "" {
		req.Phone = input.Phone.String
	} else {
		return nil, faulterr.NewFrobiddenError("phone is required")
	}

	return req, nil
}

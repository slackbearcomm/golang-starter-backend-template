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

type departmentResolver struct{ *Resolver }

// Department returns graph.DepartmentResolver implementation.
func (r *Resolver) Department() graph.DepartmentResolver { return &departmentResolver{r} }

// Organization is the resolver for the organization field.
func (r *departmentResolver) Organization(ctx context.Context, obj *dbmodels.Department) (*dbmodels.Organization, error) {
	return dataloaders.OrganizationLoaderFromContext(ctx, obj.OrgUID.String())
}

///////////////
//   Query   //
///////////////

// Departments is the resolver for the departments field.
func (r *queryResolver) Departments(ctx context.Context, search graph.SearchFilter) (*graph.DepartmentsResult, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	if search.OrgUID != nil {
		orgUID = search.OrgUID
	}

	auther, err := r.GetAutherWithPermission(ctx, models.ReadDepartment)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	filter := r.SearchFilter(search)
	output, total, err := r.services.DepartmentService.List(ctx, filter, orgUID)
	if err != nil {
		return nil, err.Error
	}
	return &graph.DepartmentsResult{Departments: output, Total: total}, nil
}

// Department is the resolver for the department field.
func (r *queryResolver) Department(ctx context.Context, id *int64, code *string) (*dbmodels.Department, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.ReadDepartment)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	if id != nil {
		obj, err := r.services.DepartmentService.GetByID(ctx, *id, orgUID)
		if err != nil {
			return nil, err.Error
		}
		return obj, nil
	}

	if code != nil {
		obj, err := r.services.DepartmentService.GetByCode(ctx, *code, orgUID)
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

// DepartmentCreate is the resolver for the departmentCreate field.
func (r *mutationResolver) DepartmentCreate(ctx context.Context, input graph.UpdateDepartment) (*dbmodels.Department, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.CreateDepartment)
	if err != nil {
		return nil, err.Error
	}

	req, reqErr := r.generateDepartmentRequest(input)
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

	obj, err := r.services.DepartmentService.Create(ctx, tx, *req)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.DepartmentObject, constants.CreateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.DepartmentObject)),
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

// DepartmentUpdate is the resolver for the departmentUpdate field.
func (r *mutationResolver) DepartmentUpdate(ctx context.Context, id int64, input graph.UpdateDepartment) (*dbmodels.Department, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateDepartment)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	req, reqErr := r.generateDepartmentRequest(input)
	if reqErr != nil {
		return nil, reqErr
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.DepartmentService.Update(ctx, tx, id, *req, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.DepartmentObject, constants.UpdateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.DepartmentObject)),
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

// DepartmentFinalize is the resolver for the departmentFinalize field.
func (r *mutationResolver) DepartmentFinalize(ctx context.Context, id int64) (*dbmodels.Department, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateDepartment)
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

	obj, err := r.services.DepartmentService.Finalize(ctx, tx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.DepartmentObject, constants.FinalizeAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.DepartmentObject)),
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

// DepartmentArchive is the resolver for the departmentArchive field.
func (r *mutationResolver) DepartmentArchive(ctx context.Context, id int64) (*dbmodels.Department, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateDepartment)
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

	obj, err := r.services.DepartmentService.Archive(ctx, tx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.DepartmentObject, constants.ArchiveAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.DepartmentObject)),
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

// DepartmentUnarchive is the resolver for the departmentUnarchive field.
func (r *mutationResolver) DepartmentUnarchive(ctx context.Context, id int64) (*dbmodels.Department, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateDepartment)
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

	obj, err := r.services.DepartmentService.Unarchive(ctx, tx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.DepartmentObject, constants.UnarchiveAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.DepartmentObject)),
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

func (r *mutationResolver) generateDepartmentRequest(input graph.UpdateDepartment) (*dbmodels.DepartmentRequest, error) {
	req := &dbmodels.DepartmentRequest{}

	if input.Name == nil || !input.Name.Valid || input.Name.String == "" {
		return nil, faulterr.NewFrobiddenError("department name is required").Error
	} else {
		req.Name = input.Name.String
	}

	return req, nil
}

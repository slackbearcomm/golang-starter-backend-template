package resolvers

import (
	"context"
	"fmt"
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/models"
	"gogql/app/models/constants"
	"gogql/app/models/dbmodels"
	"gogql/utils/faulterr"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

///////////////
//   Query   //
///////////////

func (r *queryResolver) Organizations(ctx context.Context, search graph.SearchFilter, sector *string) (*graph.OrganizationsResult, error) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		return nil, faulterr.NewUnauthorizedError(unauthorizedErr).Error
	}

	filter := r.SearchFilter(search)
	output, total, err := r.services.OrganizationService.List(ctx, filter, sector)
	if err != nil {
		return nil, err.Error
	}
	return &graph.OrganizationsResult{Organizations: output, Total: total}, nil
}

func (r *queryResolver) Organization(ctx context.Context, uid *uuid.UUID, code *string) (*dbmodels.Organization, error) {
	auther, err := r.GetAutherWithPermission(ctx, models.ReadOrganization)
	if err != nil {
		return nil, err.Error
	}

	var orgUID *uuid.UUID
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	if uid != nil {
		result, err := r.services.OrganizationService.GetByUID(ctx, *uid, orgUID)
		if err != nil {
			return nil, err.Error
		}
		return result, nil
	}

	if code != nil {
		result, err := r.services.OrganizationService.GetByCode(ctx, *code, orgUID)
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

func (r *mutationResolver) OrganizationRegister(ctx context.Context, input graph.RegisterOrganization) (*dbmodels.Organization, error) {
	// Construct Register request from input
	req := dbmodels.OrganizationRegisterRequest{}

	if input.OrgName != nil && input.OrgName.Valid && input.OrgName.String != "" {
		req.OrgName = input.OrgName.String
	} else {
		return nil, faulterr.NewFrobiddenError("org name is required").Error
	}
	if input.Website != nil && input.Website.Valid && input.Website.String != "" {
		req.Website = *input.Website
	}
	if input.Logo != nil {
		req.Logo.Name = input.Logo.Name
		req.Logo.URL = input.Logo.URL
	}
	if input.Sector != nil && input.Sector.Valid && input.Sector.String != "" {
		req.Sector = input.Sector.String
	} else {
		return nil, faulterr.NewFrobiddenError("sector is required").Error
	}

	if input.FirstName != nil && input.FirstName.Valid && input.FirstName.String != "" {
		req.FirstName = input.FirstName.String
	} else {
		return nil, faulterr.NewFrobiddenError("first name is required").Error
	}
	if input.LastName != nil && input.LastName.Valid && input.LastName.String != "" {
		req.LastName = input.LastName.String
	} else {
		return nil, faulterr.NewFrobiddenError("last name is required").Error
	}
	if input.Email != nil && input.Email.Valid && input.Email.String != "" {
		req.Email = input.Email.String
	} else {
		return nil, faulterr.NewFrobiddenError("email is required").Error
	}
	if input.Phone != nil && input.Phone.Valid && input.Phone.String != "" {
		req.Phone = input.Phone.String
	} else {
		return nil, faulterr.NewFrobiddenError("phone is required").Error
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	org, _, err := r.services.OrganizationService.Create(ctx, tx, req)
	if err != nil {
		return nil, err.Error
	}

	// commit db transaction
	if err := r.services.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err.Error
	}

	return org, nil
}

func (r *mutationResolver) OrganizationUpdate(ctx context.Context, uid uuid.UUID, input graph.UpdateOrganization) (*dbmodels.Organization, error) {
	auther, err := r.GetAutherWithPermission(ctx, models.UpdateOrganization)
	if err != nil {
		return nil, err.Error
	}

	var orgUID *uuid.UUID
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	req := dbmodels.OrganizationRequest{}

	if input.Name != nil && input.Name.Valid && input.Name.String != "" {
		req.Name = input.Name.String
	} else {
		return nil, faulterr.NewFrobiddenError("name is required").Error
	}
	if input.Sector != nil && input.Sector.Valid && input.Sector.String != "" {
		req.Sector = input.Sector.String
	} else {
		return nil, faulterr.NewFrobiddenError("sector is required").Error
	}

	if input.Website != nil && input.Website.Valid && input.Website.String != "" {
		req.Website = *input.Website
	}
	if input.Logo != nil {
		req.Logo.Name = input.Logo.Name
		req.Logo.URL = input.Logo.URL
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.OrganizationService.Update(ctx, tx, uid, req, orgUID)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.OrganizationObject, constants.UpdateAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.OrganizationObject)),
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

// OrganizationArchive is the resolver for the organizationArchive field.
func (r *mutationResolver) OrganizationArchive(ctx context.Context, uid uuid.UUID) (*dbmodels.Organization, error) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		return nil, faulterr.NewUnauthorizedError(unauthorizedErr).Error
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.OrganizationService.Archive(ctx, tx, uid)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.OrganizationObject, constants.ArchiveAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.OrganizationObject)),
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

// OrganizationUnarchive is the resolver for the organizationUnarchive field.
func (r *mutationResolver) OrganizationUnarchive(ctx context.Context, uid uuid.UUID) (*dbmodels.Organization, error) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		return nil, faulterr.NewFrobiddenError(noQueryParamsErr).Error
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	obj, err := r.services.OrganizationService.Unarchive(ctx, tx, uid)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       fmt.Sprintf("%s_%s", constants.OrganizationObject, constants.UnarchiveAction),
		ObjectID:     null.Int64From(obj.ID),
		ObjectType:   null.StringFrom(string(constants.OrganizationObject)),
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

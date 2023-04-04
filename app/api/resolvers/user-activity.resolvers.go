package resolvers

import (
	"context"
	"gogql/app/api/dataloaders"
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/middlewares"
	"gogql/app/models"
	"gogql/app/models/dbmodels"
)

type userActivityResolver struct{ *Resolver }

// UserActivity returns graph.UserActivityResolver implementation.
func (r *Resolver) UserActivity() graph.UserActivityResolver { return &userActivityResolver{r} }

// User is the resolver for the user field.
func (r *userActivityResolver) User(ctx context.Context, obj *dbmodels.UserActivity) (*dbmodels.User, error) {
	return dataloaders.UserLoaderFromContext(ctx, obj.UserID)
}

// Organization is the resolver for the organization field.
func (r *userActivityResolver) Organization(ctx context.Context, obj *dbmodels.UserActivity) (*dbmodels.Organization, error) {
	if obj.OrgUID.Valid {
		return dataloaders.OrganizationLoaderFromContext(ctx, obj.OrgUID.UUID.String())
	}
	return nil, nil
}

///////////////
//   Query   //
///////////////

// UserActivities is the resolver for the userActivities field.
func (r *queryResolver) UserActivities(ctx context.Context, search graph.SearchFilter, userID *int64) (*graph.UserActivitiesResult, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	if search.OrgUID != nil {
		orgUID = search.OrgUID
	}

	auther, err := r.GetAutherWithPermission(ctx, models.ReadUserActivity)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	filter := r.SearchFilter(search)
	output, total, err := r.services.UserActivityService.List(ctx, filter, userID, orgUID)
	if err != nil {
		return nil, err.Error
	}
	return &graph.UserActivitiesResult{UserActivities: output, Total: total}, nil
}

// UserActivity is the resolver for the userActivity field.
func (r *queryResolver) UserActivity(ctx context.Context, id int64) (*dbmodels.UserActivity, error) {
	orgUID := middlewares.GetOrgUID(ctx)
	auther, err := r.GetAutherWithPermission(ctx, models.ReadUserActivity)
	if err != nil {
		return nil, err.Error
	}
	if !auther.IsAdmin {
		orgUID = &auther.OrgUID.UUID
	}

	obj, err := r.services.UserActivityService.GetByID(ctx, id, orgUID)
	if err != nil {
		return nil, err.Error
	}
	return obj, nil
}

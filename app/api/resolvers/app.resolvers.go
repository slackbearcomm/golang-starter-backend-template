package resolvers

import (
	"context"
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/middlewares"
	"gogql/app/models"
	"gogql/app/services"
	"gogql/app/store/filestore"
	"gogql/utils/faulterr"
)

const (
	noQueryParamsErr string = "no query parameters provided"
	unauthorizedErr  string = "permission denied"
)

type Resolver struct {
	services  *services.Services
	filestore *filestore.FileStore
}

func NewResolver(s *services.Services, fs *filestore.FileStore) *Resolver {
	return &Resolver{s, fs}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *Resolver) GetAuther(ctx context.Context) (*models.Auther, *faulterr.FaultErr) {
	token := middlewares.GetSessionToken(ctx)
	if token == nil {
		return nil, faulterr.NewBadRequestError("no auth credentials provided")
	}
	return r.services.AuthService.GetAutherByToken(ctx, *token)
}

func (r *Resolver) GetAutherWithPermission(ctx context.Context, perm string) (*models.Auther, *faulterr.FaultErr) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, perm); err != nil {
		return nil, err
	}
	return auther, nil
}

func (r *Resolver) SearchFilter(search graph.SearchFilter) models.SearchFilter {
	filter := models.SearchFilter{}

	if search.SortBy != nil {
		filter.SortBy = string(*search.SortBy)
	}
	if search.SortDir != nil {
		filter.SortDir = string(*search.SortDir)
	}
	if search.Offset != nil {
		filter.Offset = search.Offset.Int
	}
	if search.Limit != nil {
		filter.Limit = search.Limit.Int
	}
	if search.Filter != nil {
		if search.Filter.String() == "All" {
			filter.IsArchived = nil
		}
		if search.Filter.String() == "Active" {
			final := true
			archived := false
			filter.IsFinal = &final
			filter.IsArchived = &archived
		}
		if search.Filter.String() == "Draft" {
			final := false
			archived := false
			filter.IsFinal = &final
			filter.IsArchived = &archived
		}
		if search.Filter.String() == "Accepted" {
			final := true
			accepted := true
			archived := false
			filter.IsFinal = &final
			filter.IsAccepted = &accepted
			filter.IsArchived = &archived
		}
		if search.Filter.String() == "Archived" {
			archived := true
			filter.IsArchived = &archived
		}
	}

	return filter
}

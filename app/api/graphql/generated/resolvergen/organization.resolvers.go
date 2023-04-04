package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.27

import (
	"context"
	"fmt"
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/models/dbmodels"

	"github.com/gofrs/uuid"
)

// OrganizationRegister is the resolver for the organizationRegister field.
func (r *mutationResolver) OrganizationRegister(ctx context.Context, input graph.RegisterOrganization) (*dbmodels.Organization, error) {
	panic(fmt.Errorf("not implemented: OrganizationRegister - organizationRegister"))
}

// OrganizationUpdate is the resolver for the organizationUpdate field.
func (r *mutationResolver) OrganizationUpdate(ctx context.Context, uid uuid.UUID, input graph.UpdateOrganization) (*dbmodels.Organization, error) {
	panic(fmt.Errorf("not implemented: OrganizationUpdate - organizationUpdate"))
}

// OrganizationArchive is the resolver for the organizationArchive field.
func (r *mutationResolver) OrganizationArchive(ctx context.Context, uid uuid.UUID) (*dbmodels.Organization, error) {
	panic(fmt.Errorf("not implemented: OrganizationArchive - organizationArchive"))
}

// OrganizationUnarchive is the resolver for the organizationUnarchive field.
func (r *mutationResolver) OrganizationUnarchive(ctx context.Context, uid uuid.UUID) (*dbmodels.Organization, error) {
	panic(fmt.Errorf("not implemented: OrganizationUnarchive - organizationUnarchive"))
}

// Organizations is the resolver for the organizations field.
func (r *queryResolver) Organizations(ctx context.Context, search graph.SearchFilter, sector *string) (*graph.OrganizationsResult, error) {
	panic(fmt.Errorf("not implemented: Organizations - organizations"))
}

// Organization is the resolver for the organization field.
func (r *queryResolver) Organization(ctx context.Context, uid *uuid.UUID, code *string) (*dbmodels.Organization, error) {
	panic(fmt.Errorf("not implemented: Organization - organization"))
}
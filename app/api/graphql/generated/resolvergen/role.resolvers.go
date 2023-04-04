package resolvergen

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.27

import (
	"context"
	"fmt"
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/models/dbmodels"
)

// RoleCreate is the resolver for the roleCreate field.
func (r *mutationResolver) RoleCreate(ctx context.Context, input graph.UpdateRole) (*dbmodels.Role, error) {
	panic(fmt.Errorf("not implemented: RoleCreate - roleCreate"))
}

// RoleUpdate is the resolver for the roleUpdate field.
func (r *mutationResolver) RoleUpdate(ctx context.Context, id int64, input graph.UpdateRole) (*dbmodels.Role, error) {
	panic(fmt.Errorf("not implemented: RoleUpdate - roleUpdate"))
}

// RoleFinalize is the resolver for the roleFinalize field.
func (r *mutationResolver) RoleFinalize(ctx context.Context, id int64) (*dbmodels.Role, error) {
	panic(fmt.Errorf("not implemented: RoleFinalize - roleFinalize"))
}

// RoleArchive is the resolver for the roleArchive field.
func (r *mutationResolver) RoleArchive(ctx context.Context, id int64) (*dbmodels.Role, error) {
	panic(fmt.Errorf("not implemented: RoleArchive - roleArchive"))
}

// RoleUnarchive is the resolver for the roleUnarchive field.
func (r *mutationResolver) RoleUnarchive(ctx context.Context, id int64) (*dbmodels.Role, error) {
	panic(fmt.Errorf("not implemented: RoleUnarchive - roleUnarchive"))
}

// Roles is the resolver for the roles field.
func (r *queryResolver) Roles(ctx context.Context, search graph.SearchFilter, deptID *int64) (*graph.RolesResult, error) {
	panic(fmt.Errorf("not implemented: Roles - roles"))
}

// Role is the resolver for the role field.
func (r *queryResolver) Role(ctx context.Context, id *int64, code *string) (*dbmodels.Role, error) {
	panic(fmt.Errorf("not implemented: Role - role"))
}

// Organization is the resolver for the organization field.
func (r *roleResolver) Organization(ctx context.Context, obj *dbmodels.Role) (*dbmodels.Organization, error) {
	panic(fmt.Errorf("not implemented: Organization - organization"))
}

// Department is the resolver for the department field.
func (r *roleResolver) Department(ctx context.Context, obj *dbmodels.Role) (*dbmodels.Department, error) {
	panic(fmt.Errorf("not implemented: Department - department"))
}

// Role returns graph.RoleResolver implementation.
func (r *Resolver) Role() graph.RoleResolver { return &roleResolver{r} }

type roleResolver struct{ *Resolver }
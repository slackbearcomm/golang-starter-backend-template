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

// Organization is the resolver for the organization field.
func (r *departmentResolver) Organization(ctx context.Context, obj *dbmodels.Department) (*dbmodels.Organization, error) {
	panic(fmt.Errorf("not implemented: Organization - organization"))
}

// DepartmentCreate is the resolver for the departmentCreate field.
func (r *mutationResolver) DepartmentCreate(ctx context.Context, input graph.UpdateDepartment) (*dbmodels.Department, error) {
	panic(fmt.Errorf("not implemented: DepartmentCreate - departmentCreate"))
}

// DepartmentUpdate is the resolver for the departmentUpdate field.
func (r *mutationResolver) DepartmentUpdate(ctx context.Context, id int64, input graph.UpdateDepartment) (*dbmodels.Department, error) {
	panic(fmt.Errorf("not implemented: DepartmentUpdate - departmentUpdate"))
}

// DepartmentFinalize is the resolver for the departmentFinalize field.
func (r *mutationResolver) DepartmentFinalize(ctx context.Context, id int64) (*dbmodels.Department, error) {
	panic(fmt.Errorf("not implemented: DepartmentFinalize - departmentFinalize"))
}

// DepartmentArchive is the resolver for the departmentArchive field.
func (r *mutationResolver) DepartmentArchive(ctx context.Context, id int64) (*dbmodels.Department, error) {
	panic(fmt.Errorf("not implemented: DepartmentArchive - departmentArchive"))
}

// DepartmentUnarchive is the resolver for the departmentUnarchive field.
func (r *mutationResolver) DepartmentUnarchive(ctx context.Context, id int64) (*dbmodels.Department, error) {
	panic(fmt.Errorf("not implemented: DepartmentUnarchive - departmentUnarchive"))
}

// Departments is the resolver for the departments field.
func (r *queryResolver) Departments(ctx context.Context, search graph.SearchFilter) (*graph.DepartmentsResult, error) {
	panic(fmt.Errorf("not implemented: Departments - departments"))
}

// Department is the resolver for the department field.
func (r *queryResolver) Department(ctx context.Context, id *int64, code *string) (*dbmodels.Department, error) {
	panic(fmt.Errorf("not implemented: Department - department"))
}

// Department returns graph.DepartmentResolver implementation.
func (r *Resolver) Department() graph.DepartmentResolver { return &departmentResolver{r} }

type departmentResolver struct{ *Resolver }

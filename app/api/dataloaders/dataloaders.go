package dataloaders

import (
	"context"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"net/http"
	"time"
)

// ContextKey holds a custom String func for uniqueness
type ContextKey string

func (k ContextKey) String() string {
	return "dataloader_" + string(k)
}

// OrganizationLoaderKey declares a statically typed key for context reference in other packages
const OrganizationLoaderKey ContextKey = "organization_loader"

// DepartmentLoaderKey declares a statically typed key for context reference in other packages
const DepartmentLoaderKey ContextKey = "department_loader"

// RoleLoaderKey declares a statically typed key for context reference in other packages
const RoleLoaderKey ContextKey = "role_loader"

// UserLoaderKey declares a statically typed key for context reference in other packages
const UserLoaderKey ContextKey = "user_loader"

// OrganizationLoaderFromContext runs the dataloader inside the context
func OrganizationLoaderFromContext(ctx context.Context, uid string) (*dbmodels.Organization, error) {
	return ctx.Value(OrganizationLoaderKey).(*OrganizationLoader).Load(uid)
}

// DepartmentLoaderFromContext runs the dataloader inside the context
func DepartmentLoaderFromContext(ctx context.Context, id int64) (*dbmodels.Department, error) {
	return ctx.Value(DepartmentLoaderKey).(*DepartmentLoader).Load(id)
}

// RoleLoaderFromContext runs the dataloader inside the context
func RoleLoaderFromContext(ctx context.Context, id int64) (*dbmodels.Role, error) {
	return ctx.Value(RoleLoaderKey).(*RoleLoader).Load(id)
}

// UserLoaderFromContext runs the dataloader inside the context
func UserLoaderFromContext(ctx context.Context, id int64) (*dbmodels.User, error) {
	return ctx.Value(UserLoaderKey).(*UserLoader).Load(id)
}

// WithDataloaders returns a new context that contains dataloaders
func WithDataloaders(
	ctx context.Context,
	dbstore *dbstore.DBStore,
) context.Context {
	organizationLoader := NewOrganizationLoader(
		OrganizationLoaderConfig{
			Fetch: func(uids []string) ([]*dbmodels.Organization, []error) {
				data, err := dbstore.OrganizationStore.GetManyByUIDs(ctx, uids)
				if err != nil {
					return nil, []error{err.Error}
				}

				// make result and ids of the same order
				slice := make(map[interface{}]*dbmodels.Organization, len(data))
				for _, e := range data {
					slice[e.UID.String()] = e
				}

				result := make([]*dbmodels.Organization, len(uids))
				for i, key := range uids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	departmentLoader := NewDepartmentLoader(
		DepartmentLoaderConfig{
			Fetch: func(ids []int64) ([]*dbmodels.Department, []error) {
				data, err := dbstore.DepartmentStore.GetManyByIDs(ctx, ids)
				if err != nil {
					return nil, []error{err.Error}
				}

				// make result and ids of the same order
				slice := make(map[interface{}]*dbmodels.Department, len(data))
				for _, e := range data {
					slice[e.ID] = e
				}

				result := make([]*dbmodels.Department, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	roleLoader := NewRoleLoader(
		RoleLoaderConfig{
			Fetch: func(ids []int64) ([]*dbmodels.Role, []error) {
				data, err := dbstore.RoleStore.GetManyByIDs(ctx, ids)
				if err != nil {
					return nil, []error{err.Error}
				}

				// make result and ids of the same order
				slice := make(map[interface{}]*dbmodels.Role, len(data))
				for _, e := range data {
					slice[e.ID] = e
				}

				result := make([]*dbmodels.Role, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	userLoader := NewUserLoader(
		UserLoaderConfig{
			Fetch: func(ids []int64) ([]*dbmodels.User, []error) {
				data, err := dbstore.UserStore.GetManyByIDs(ctx, ids)
				if err != nil {
					return nil, []error{err.Error}
				}

				// make result and ids of the same order
				slice := make(map[interface{}]*dbmodels.User, len(data))
				for _, e := range data {
					slice[e.ID] = e
				}

				result := make([]*dbmodels.User, len(ids))
				for i, key := range ids {
					result[i] = slice[key]
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		},
	)

	ctx = context.WithValue(ctx, OrganizationLoaderKey, organizationLoader)
	ctx = context.WithValue(ctx, DepartmentLoaderKey, departmentLoader)
	ctx = context.WithValue(ctx, RoleLoaderKey, roleLoader)
	ctx = context.WithValue(ctx, UserLoaderKey, userLoader)
	return ctx
}

// DataloaderMiddleware runs before each API call and loads the dataloaders into context
func DataloaderMiddleware(
	dbstore *dbstore.DBStore,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(WithDataloaders(r.Context(), dbstore))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

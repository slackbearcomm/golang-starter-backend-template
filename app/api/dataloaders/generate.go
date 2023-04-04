//go:generate go run github.com/vektah/dataloaden OrganizationLoader string *gogql/app/models/dbmodels.Organization
//go:generate go run github.com/vektah/dataloaden DepartmentLoader int64 *gogql/app/models/dbmodels.Department
//go:generate go run github.com/vektah/dataloaden RoleLoader int64 *gogql/app/models/dbmodels.Role
//go:generate go run github.com/vektah/dataloaden UserLoader int64 *gogql/app/models/dbmodels.User

package dataloaders

package services

import (
	"gogql/app/master"
	"gogql/app/services/authservice"
	"gogql/app/services/orgservice"
	"gogql/app/services/settingservice"
	"gogql/app/store/dbstore"
)

type Services struct {
	// settings
	DBTX *settingservice.DBTX

	// authentication
	AuthService *authservice.AuthService

	// companies
	OrganizationService *orgservice.OrganizationService
	DepartmentService   *orgservice.DepartmentService
	RoleService         *orgservice.RoleService
	UserService         *orgservice.UserService
	UserActivityService *orgservice.UserActivityService
}

func NewService(dbs *dbstore.DBStore, master *master.Master) *Services {
	return &Services{
		// settings
		settingservice.NewDBTX(dbs),

		// authentication
		authservice.NewAuthService(dbs, master),

		// companies
		orgservice.NewOrganizationService(dbs, master),
		orgservice.NewDepartmentService(dbs, master),
		orgservice.NewRoleService(dbs, master),
		orgservice.NewUserService(dbs, master),
		orgservice.NewUserActivityService(dbs, master),
	}
}

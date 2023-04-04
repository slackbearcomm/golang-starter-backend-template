package master

import (
	"gogql/app/master/orgmaster"
	"gogql/app/store/dbstore"
)

type Master struct {
	// companies
	OrganizationMaster *orgmaster.OrganizationMaster
	DepartmentMaster   *orgmaster.DepartmentMaster
	RoleMaster         *orgmaster.RoleMaster
	UserMaster         *orgmaster.UserMaster
	OTPSessionMaster   *orgmaster.OTPSessionMaster
	AuthSessionMaster  *orgmaster.AuthSessionMaster
	UserActivityMaster *orgmaster.UserActivityMaster
}

func NewMaster(dbStore *dbstore.DBStore) *Master {
	return &Master{
		// companies
		orgmaster.NewOrganizationMaster(dbStore),
		orgmaster.NewDepartmentMaster(dbStore),
		orgmaster.NewRoleMaster(dbStore),
		orgmaster.NewUserMaster(dbStore),
		orgmaster.NewOTPSessionMaster(dbStore),
		orgmaster.NewAuthSessionMaster(dbStore),
		orgmaster.NewUserActivityMaster(dbStore),
	}
}

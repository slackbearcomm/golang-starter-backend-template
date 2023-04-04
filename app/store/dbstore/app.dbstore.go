package dbstore

import (
	"gogql/app/store/dbstore/orgstore"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStore struct {
	// settings
	DBTX *DBTX

	// comapnies
	OrganizationStore *orgstore.OrganizationStore
	DepartmentStore   *orgstore.DepartmentStore
	RoleStore         *orgstore.RoleStore
	UserStore         *orgstore.UserStore
	OTPSessionStore   *orgstore.OTPSessionStore
	AuthSessionStore  *orgstore.AuthSessionStore
	UserActivityStore *orgstore.UserActivityStore
}

func NewDBStore(conn *pgxpool.Pool) *DBStore {
	return &DBStore{
		// settings
		NewDBTX(conn),

		// comapnies
		orgstore.NewOrganizationStore(conn),
		orgstore.NewDepartmentStore(conn),
		orgstore.NewRoleStore(conn),
		orgstore.NewUserStore(conn),
		orgstore.NewOTPSessionStore(conn),
		orgstore.NewAuthSessionStore(conn),
		orgstore.NewUserActivityStore(conn),
	}
}

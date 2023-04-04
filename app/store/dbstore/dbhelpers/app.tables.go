package dbhelpers

type dbTable string

const (
	// comapnies
	OrganizationsTable  dbTable = "organizations"
	DepartmentsTable    dbTable = "departments"
	RolesTable          dbTable = "roles"
	UsersTable          dbTable = "users"
	UserActivitiesTable dbTable = "user_activities"
)

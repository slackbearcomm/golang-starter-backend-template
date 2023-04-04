package models

const (
	UploadFile string = "UPLOAD_FILE"

	// companies
	CreateOrganization string = "CREATE_ORGANIZATION"
	ReadOrganization   string = "READ_ORGANIZATION"
	UpdateOrganization string = "UPDATE_ORGANIZATION"
	DeleteOrganization string = "DELETE_ORGANIZATION"
	CreateDepartment   string = "CREATE_DEPARTMENT"
	ReadDepartment     string = "READ_DEPARTMENT"
	UpdateDepartment   string = "UPDATE_DEPARTMENT"
	DeleteDepartment   string = "DELETE_DEPARTMENT"
	CreateRole         string = "CREATE_ROLE"
	ReadRole           string = "READ_ROLE"
	UpdateRole         string = "UPDATE_ROLE"
	DeleteRole         string = "DELETE_ROLE"
	CreateUser         string = "CREATE_USER"
	ReadUser           string = "READ_USER"
	UpdateUser         string = "UPDATE_USER"
	DeleteUser         string = "DELETE_USER"
	CreateContact      string = "CREATE_CONTRACT"
	ReadContact        string = "READ_CONTRACT"
	UpdateContact      string = "UPDATE_CONTRACT"
	DeleteContact      string = "DELETE_CONTRACT"
	ReadUserActivity   string = "READ_USER_ACTIVITY"
)

func ListPermissions() []string {
	return []string{
		UploadFile,
		CreateOrganization,
		ReadOrganization,
		UpdateOrganization,
		DeleteOrganization,
		CreateDepartment,
		ReadDepartment,
		UpdateDepartment,
		DeleteDepartment,
		CreateRole,
		ReadRole,
		UpdateRole,
		DeleteRole,
		CreateUser,
		ReadUser,
		UpdateUser,
		DeleteUser,
	}
}

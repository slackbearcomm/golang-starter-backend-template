package constants

type ObjectType string

const (
	LoginAction     string = "LOGIN"
	CreateAction    string = "CREATE"
	UpdateAction    string = "UPDATE"
	FinalizeAction  string = "FINALIZE"
	AcceptAction    string = "ACCEPT"
	DeclineAction   string = "DECLINE"
	ArchiveAction   string = "ARCHIVE"
	UnarchiveAction string = "UNARCHIVE"
)

const (
	SelfObject   ObjectType = "SELF"
	AutherObject ObjectType = "AUTHER"
	FileObject   ObjectType = "FILE"

	// Company
	OrganizationObject ObjectType = "ORGANIZATION"
	DepartmentObject   ObjectType = "DEPARTMENT"
	RoleObject         ObjectType = "ROLE"
	UserObject         ObjectType = "USER"
	ContactObject      ObjectType = "CONTACT"
)

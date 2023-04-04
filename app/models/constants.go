package models

const (
	DeptManagementStr string = "Management"
	RoleManagementStr string = "Management"
	RoleStaffStr      string = "Staff"
)

type SortByOption string
type SortDir string

const (
	SortByOptionDateCreated  SortByOption = "DateCreated"
	SortByOptionDateUpdated  SortByOption = "DateUpdated"
	SortByOptionAlphabetical SortByOption = "Alphabetical"
)

const (
	SortDirAscending  SortDir = "Ascending"
	SortDirDescending SortDir = "Descending"
)

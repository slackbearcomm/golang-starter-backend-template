package helpers

import "gogql/app/models"

func GetStaffPerm() []string {
	return []string{
		models.ReadOrganization,
		models.ReadRole,
		models.ReadUser,
	}
}

package helpers

import (
	"fmt"
	"gogql/app/models/constants"
	"gogql/utils/faulterr"

	"github.com/gofrs/uuid"
)

// GenerateUID generates a uuid v4
func GenerateUID() (*uuid.UUID, *faulterr.FaultErr) {
	uid, uidErr := uuid.NewV4()
	if uidErr != nil {
		return nil, faulterr.NewInternalServerError(uidErr.Error())
	}

	return &uid, nil
}

// NullUUIDFromUUID generates NullUUID from a valid UUID
func NullUUIDFromUUID(uid uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: uid, Valid: true}
}

// GenerateCode generates code value
func GenerateCode(obj constants.ObjectType, orgCode string, count int64) string {
	var orgcode string
	// trim org code
	if orgCode != "" {
		orgcode = string([]rune(orgCode)[3:])
	}

	// increase count by 1
	count = count + 1
	switch obj {
	// company
	case constants.OrganizationObject:
		return fmt.Sprintf("ORG%03d", count)
	case constants.DepartmentObject:
		return fmt.Sprintf("DEPT-%s-%03d", orgcode, count)
	case constants.RoleObject:
		return fmt.Sprintf("ROLE-%s-%03d", orgcode, count)
	case constants.ContactObject:
		return fmt.Sprintf("CONT-%s-%03d", orgcode, count)
	// default
	default:
		return fmt.Sprintf("-%s-%03d", orgcode, count)
	}
}

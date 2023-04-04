package orgseed

import (
	"context"
	"fmt"
	"gogql/app/helpers"
	"gogql/app/master"
	"gogql/app/models"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/seed/seedconst"
	"gogql/utils/faulterr"
	"gogql/utils/logger"

	"github.com/jackc/pgx/v5"
	"github.com/volatiletech/null"
)

const (
	ABL            string = "ABL"
	GREENDAY       string = "Greenday"
	AV             string = "AV"
	WAREPORT       string = "Wareport"
	CAPITAL_BRIDGE string = "Capital Bridge"
)

func InsertOrganizations(ctx context.Context, tx pgx.Tx, d *dbstore.DBStore, m *master.Master) ([]*dbmodels.Organization, *faulterr.FaultErr) {
	orgRequests := []dbmodels.OrganizationRequest{}
	for _, r := range orgRegsiterReqs {
		r := dbmodels.OrganizationRequest{
			Name:    r.OrgName,
			Website: r.Website,
			Logo:    r.Logo,
			Sector:  r.Sector,
		}
		orgRequests = append(orgRequests, r)
	}

	organizations, err := m.OrganizationMaster.BulkCreate(ctx, tx, orgRequests)
	if err != nil {
		return nil, err
	}

	logger.Success(fmt.Sprintf("%v organizations added to the database", len(organizations)))

	// insert departments
	departments := []*dbmodels.Department{}
	for i := range organizations {
		deptReq := dbmodels.DepartmentRequest{
			OrgUID:  organizations[i].UID,
			Name:    models.DeptManagementStr,
			IsFinal: true,
		}

		dept, err := m.DepartmentMaster.CreateOne(ctx, tx, deptReq, *organizations[i])
		if err != nil {
			return nil, err
		}

		departments = append(departments, dept)
	}

	// insert roles
	roles := []*dbmodels.Role{}
	for i, org := range organizations {
		roleReq := dbmodels.RoleRequest{
			OrgUID:       org.UID,
			DepartmentID: departments[i].ID,
			Name:         models.RoleManagementStr,
			Permissions:  nil,
			IsManagement: true,
			IsFinal:      true,
		}
		role, err := m.RoleMaster.CreateOne(ctx, tx, roleReq, *organizations[i])
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	// insert org admins
	for i, e := range orgRegsiterReqs {
		userReq := dbmodels.UserRequest{
			FirstName: e.FirstName,
			LastName:  e.LastName,
			Email:     e.Email,
			Phone:     e.Phone,
			OrgUID:    helpers.NullUUIDFromUUID(roles[i].OrgUID),
			RoleID:    null.Int64From(roles[i].ID),
			IsFinal:   true,
		}

		_, err = m.UserMaster.CreateOne(ctx, tx, userReq)
		if err != nil {
			return nil, err
		}
	}

	logger.Success(fmt.Sprintf(
		"%v departments, %v roles, %v users added to the database",
		len(departments),
		len(roles),
		len(orgRegsiterReqs),
	))

	return organizations, nil
}

// Seed data
var orgRegsiterReqs = []dbmodels.OrganizationRegisterRequest{
	{
		OrgName:   seedconst.ABL,
		Website:   null.StringFrom("www.abl.com"),
		FirstName: "Abhinandan",
		LastName:  "Darbey",
		Email:     "abl@example.com",
		Phone:     "1234567890",
	},
	{
		OrgName:   seedconst.GREENDAY,
		Website:   null.StringFrom("www.greenday.com"),
		FirstName: "Saurabh",
		LastName:  "Chakrabarty",
		Email:     "greenday@example.com",
		Phone:     "1234567891",
	},
	{
		OrgName:   seedconst.AV,
		Website:   null.StringFrom("www.av.com"),
		FirstName: "Ajay",
		LastName:  "Vishwakarma",
		Email:     "av@example.com",
		Phone:     "1234567892",
	},
	{
		OrgName:   seedconst.WAREPORT,
		Website:   null.StringFrom("www.wareport.com"),
		FirstName: "Vedant",
		LastName:  "Pande",
		Email:     "wareport@example.com",
		Phone:     "1234567893",
	},
	{
		OrgName:   seedconst.CAPITAL_BRIDGE,
		Website:   null.StringFrom("www.capitalbridge.com"),
		FirstName: "Abhimanyu",
		LastName:  "Darbey",
		Email:     "capitalbridge@example.com",
		Phone:     "1234567894",
	},
}

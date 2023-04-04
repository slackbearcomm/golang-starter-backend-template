package orgseed

import (
	"context"
	"fmt"
	"gogql/app/master"
	"gogql/app/models/dbmodels"
	"gogql/utils/faulterr"
	"gogql/utils/logger"

	"github.com/jackc/pgx/v5"
)

func InsertAdmin(tx pgx.Tx, m *master.Master) (*dbmodels.User, *faulterr.FaultErr) {
	admins := []dbmodels.UserRequest{
		{
			FirstName: "Super",
			LastName:  "Admin",
			Email:     "superadmin@example.com",
			Phone:     "9000090000",
			IsAdmin:   true,
		},
	}

	superAdmins := []*dbmodels.User{}

	for _, r := range admins {
		ctx := context.Background()
		obj, err := m.UserMaster.CreateOne(ctx, tx, r)
		if err != nil {
			return nil, err
		}

		superAdmins = append(superAdmins, obj)
	}

	logger.Success(fmt.Sprintf("%v super admins added to the database", len(admins)))

	return superAdmins[0], nil
}

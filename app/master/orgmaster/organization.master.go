package orgmaster

import (
	"context"
	"gogql/app/helpers"
	"gogql/app/models/constants"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type OrganizationMaster struct {
	dbstore *dbstore.DBStore
}

func NewOrganizationMaster(dbstore *dbstore.DBStore) *OrganizationMaster {
	return &OrganizationMaster{dbstore}
}

func (m *OrganizationMaster) VerifyOrganizationExists(ctx context.Context, uid uuid.UUID) (*dbmodels.Organization, *faulterr.FaultErr) {
	// verify if org exists
	org, err := m.dbstore.OrganizationStore.GetByUID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if org.IsArchived {
		return nil, faulterr.NewNotFoundError("organization does not exists")
	}
	return org, nil
}

func (m *OrganizationMaster) CreateOne(ctx context.Context, tx pgx.Tx, req dbmodels.OrganizationRequest) (*dbmodels.Organization, *faulterr.FaultErr) {
	requests := []dbmodels.OrganizationRequest{}
	requests = append(requests, req)

	result, err := m.BulkCreate(ctx, tx, requests)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func (m *OrganizationMaster) BulkCreate(ctx context.Context, tx pgx.Tx, requests []dbmodels.OrganizationRequest) ([]*dbmodels.Organization, *faulterr.FaultErr) {
	result := []*dbmodels.Organization{}

	// Get count
	count, err := m.dbstore.OrganizationStore.GetCount(ctx)
	if err != nil {
		return nil, err
	}

	for i := range requests {
		// validate request
		if err := m.validate(requests[i]); err != nil {
			return nil, err
		}

		// generate uid
		uid, err := helpers.GenerateUID()
		if err != nil {
			return nil, err
		}

		// construct arguments
		arg := m.construct(requests[i])
		arg.UID = *uid
		arg.Code = helpers.GenerateCode(constants.OrganizationObject, "", count)
		if arg.Status == "" {
			arg.Status = constants.StatusActive
		}

		// insert into db
		obj, err := m.dbstore.OrganizationStore.Insert(ctx, tx, *arg)
		if err != nil {
			return nil, err
		}
		result = append(result, obj)
		count++
	}

	return result, nil
}

func (m *OrganizationMaster) construct(r dbmodels.OrganizationRequest) *dbmodels.Organization {
	return &dbmodels.Organization{
		Name:       r.Name,
		Website:    r.Website,
		Logo:       r.Logo,
		Sector:     r.Sector,
		IsFinal:    true,
		IsArchived: false,
	}
}

func (m *OrganizationMaster) validate(r dbmodels.OrganizationRequest) *faulterr.FaultErr {
	if r.Name == "" {
		return faulterr.NewBadRequestError("Organization Name is required")
	}
	return nil
}

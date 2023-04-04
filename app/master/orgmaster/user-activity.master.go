package orgmaster

import (
	"context"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/jackc/pgx/v5"
)

type UserActivityMaster struct {
	dbstore *dbstore.DBStore
}

func NewUserActivityMaster(s *dbstore.DBStore) *UserActivityMaster {
	return &UserActivityMaster{s}
}

func (m *UserActivityMaster) Create(ctx context.Context, tx pgx.Tx, req dbmodels.UserActivityRequest) (*dbmodels.UserActivity, *faulterr.FaultErr) {
	obj, err := m.construct(req)
	if err != nil {
		return nil, err
	}
	return m.dbstore.UserActivityStore.Insert(ctx, tx, *obj)
}

func (m *UserActivityMaster) construct(req dbmodels.UserActivityRequest) (*dbmodels.UserActivity, *faulterr.FaultErr) {
	obj := &dbmodels.UserActivity{
		ObjectID:     req.ObjectID,
		ObjectType:   req.ObjectType,
		OrgUID:       req.OrgUID,
		SessionToken: req.SessionToken,
	}

	if req.UserID > 0 {
		obj.UserID = req.UserID
	} else {
		return nil, faulterr.NewBadRequestError("user id is required")
	}
	if req.Action != "" {
		obj.Action = req.Action
	} else {
		return nil, faulterr.NewBadRequestError("action is required")
	}
	return obj, nil
}

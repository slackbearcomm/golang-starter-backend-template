package orgservice

import (
	"context"
	"gogql/app/master"
	"gogql/app/models"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type UserActivityService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

var _ UserActivityServiceInterface = &UserActivityService{}

type UserActivityServiceInterface interface {
	List(ctx context.Context, filter models.SearchFilter, userID *int64, orgUID *uuid.UUID) ([]dbmodels.UserActivity, int, *faulterr.FaultErr)
	GetByID(ctx context.Context, id int64, orgUID *uuid.UUID) (*dbmodels.UserActivity, *faulterr.FaultErr)
}

func NewUserActivityService(s *dbstore.DBStore, m *master.Master) *UserActivityService {
	return &UserActivityService{s, m}
}

// List gets all user activities
func (s *UserActivityService) List(ctx context.Context, filter models.SearchFilter, userID *int64, orgUID *uuid.UUID) ([]dbmodels.UserActivity, int, *faulterr.FaultErr) {
	return s.dbstore.UserActivityStore.List(ctx, filter, userID, orgUID)
}

func (s *UserActivityService) GetByID(ctx context.Context, id int64, orgUID *uuid.UUID) (*dbmodels.UserActivity, *faulterr.FaultErr) {
	activity, err := s.dbstore.UserActivityStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if orgUID != nil && activity.OrgUID.Valid && activity.OrgUID.UUID != *orgUID {
		return nil, faulterr.NewUnauthorizedError("permission denied")
	}

	return activity, nil
}

func (s *UserActivityService) Create(ctx context.Context, tx pgx.Tx, req dbmodels.UserActivityRequest) (*dbmodels.UserActivity, *faulterr.FaultErr) {
	return s.master.UserActivityMaster.Create(ctx, tx, req)
}

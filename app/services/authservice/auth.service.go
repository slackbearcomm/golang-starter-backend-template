package authservice

import (
	"context"
	"fmt"
	"gogql/app/helpers"
	"gogql/app/master"
	"gogql/app/models"
	"gogql/app/models/dbmodels"
	"gogql/app/store/dbstore"
	"gogql/utils/faulterr"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	dbstore *dbstore.DBStore
	master  *master.Master
}

func NewAuthService(dbstore *dbstore.DBStore, master *master.Master) *AuthService {
	return &AuthService{dbstore, master}
}

// GetOTP validates validates user by email and return an otp string
func (s *AuthService) GetOTP(ctx context.Context, tx pgx.Tx, req *models.OTPRequest) (*string, *faulterr.FaultErr) {
	user := &dbmodels.User{}
	if req.Email.Valid && req.Email.String != "" {
		obj, err := s.dbstore.UserStore.GetByEmail(ctx, req.Email.String)
		if err != nil {
			if err.Status == http.StatusNotFound {
				return nil, faulterr.NewFrobiddenError("no user found with given email")
			}
			return nil, err
		}
		user = obj
	}

	if req.Phone.Valid && req.Email.String != "" {
		obj, err := s.dbstore.UserStore.GetByPhone(ctx, req.Phone.String)
		if err != nil {
			return nil, err
		}
		user = obj
	}

	// generate OTP
	otp, err := s.master.OTPSessionMaster.Create(ctx, tx, user.ID)
	if err != nil {
		return nil, err
	}

	return &otp.Token, nil
}

// Login validates password and returns user
func (s *AuthService) Login(ctx context.Context, tx pgx.Tx, req *models.LoginRequest) (*models.Auther, *faulterr.FaultErr) {
	user := &dbmodels.User{}
	if req.Email.Valid && req.Email.String != "" {
		obj, err := s.dbstore.UserStore.GetByEmail(ctx, req.Email.String)
		if err != nil {
			if err.Status == http.StatusNotFound {
				return nil, faulterr.NewFrobiddenError("no user found with given email")
			}
			return nil, err
		}
		user = obj
	}
	if req.Phone.Valid && req.Email.String != "" {
		obj, err := s.dbstore.UserStore.GetByPhone(ctx, req.Phone.String)
		if err != nil {
			if err.Status == http.StatusNotFound {
				return nil, faulterr.NewFrobiddenError("no user found with given phone")
			}
			return nil, err
		}
		user = obj
	}

	// get otp from db and validate
	otp, err := s.dbstore.OTPSessionStore.GetByToken(ctx, req.OTP)
	if err != nil {
		if err.Status == http.StatusNotFound {
			return nil, faulterr.NewFrobiddenError("otp is invalid")
		}
		return nil, err
	}
	if !otp.IsValid {
		return nil, faulterr.NewUnauthorizedError("otp is not valid")
	}
	if err := helpers.ValidateTokenExpiry(otp.ExpiresAt); err != nil {
		return nil, err
	}

	// update otp validity
	otp.IsValid = false
	if err := s.dbstore.OTPSessionStore.Update(ctx, tx, otp); err != nil {
		return nil, err
	}

	// generate auth session token
	authSession, err := s.master.AuthSessionMaster.Create(ctx, tx, user.ID)
	if err != nil {
		return nil, err
	}
	return s.getAuther(user, authSession.Token), nil
}

func (s *AuthService) GetAutherByToken(ctx context.Context, token uuid.UUID) (*models.Auther, *faulterr.FaultErr) {
	authSession, err := s.dbstore.AuthSessionStore.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if !authSession.IsValid {
		return nil, faulterr.NewUnauthorizedError("token is not valid")
	}
	if err := helpers.ValidateTokenExpiry(authSession.ExpiresAt); err != nil {
		return nil, err
	}

	// get user
	user, err := s.dbstore.UserStore.GetByID(ctx, authSession.UserID)
	if err != nil {
		return nil, err
	}
	return s.getAuther(user, token), nil
}

// Helpers

func (s *AuthService) getAuther(u *dbmodels.User, token uuid.UUID) *models.Auther {
	name := fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	return &models.Auther{
		ID:           u.ID,
		Name:         name,
		IsAdmin:      u.IsAdmin,
		OrgUID:       u.OrgUID,
		RoleID:       u.RoleID,
		SessionToken: token,
	}
}

// GrantPermission verifies the member's permission and returns unauthorized error if not permitted
func (s *AuthService) GrantPermission(ctx context.Context, auther *models.Auther, perm string) *faulterr.FaultErr {
	errMsg := "permission denied"

	if !auther.IsAdmin {
		// get user role and permissions
		role, err := s.dbstore.RoleStore.GetByID(ctx, auther.RoleID.Int64)
		if err != nil {
			return err
		}

		if !role.IsManagement {
			for _, permission := range role.Permissions {
				if permission == perm {
					return nil
				}
			}
			return faulterr.NewUnauthorizedError(errMsg)
		}
		return nil
	}
	return nil
}

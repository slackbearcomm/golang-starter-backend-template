package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/models"
	"gogql/app/models/constants"
	"gogql/app/models/dbmodels"
	"gogql/utils/faulterr"

	"github.com/volatiletech/null"
)

// Auther is the resolver for the auther field.
func (r *queryResolver) Auther(ctx context.Context) (*models.Auther, error) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	return auther, nil
}

// GenerateOtp is the resolver for the generateOtp field.
func (r *mutationResolver) GenerateOtp(ctx context.Context, input *graph.OTPRequest) (*string, error) {
	req := &models.OTPRequest{}

	if input.Email != nil && input.Email.String != "" {
		req.Email = *input.Email
	}
	if input.Phone != nil && input.Phone.String != "" {
		req.Phone = *input.Phone
	}
	if !req.Email.Valid && !req.Phone.Valid {
		return nil, faulterr.NewFrobiddenError("email or phone is required").Error
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	otp, err := r.services.AuthService.GetOTP(ctx, tx, req)
	if err != nil {
		return nil, err.Error
	}

	// commit db transaction
	if err := r.services.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err.Error
	}

	return otp, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input graph.LoginRequest) (*models.Auther, error) {
	req := &models.LoginRequest{}

	if input.Email != nil && input.Email.String != "" {
		req.Email = *input.Email
	}
	if input.Phone != nil && input.Phone.String != "" {
		req.Phone = *input.Phone
	}
	if !req.Email.Valid && !req.Phone.Valid {
		return nil, faulterr.NewFrobiddenError("email or phone is required").Error
	}

	if input.Otp != nil && input.Otp.String != "" {
		req.OTP = input.Otp.String
	} else {
		return nil, faulterr.NewFrobiddenError("otp is required").Error
	}

	// start db transaction
	tx, err := r.services.DBTX.BeginTx(ctx)
	if err != nil {
		return nil, err.Error
	}
	defer r.services.DBTX.RollbackTx(ctx, tx)

	auther, err := r.services.AuthService.Login(ctx, tx, req)
	if err != nil {
		return nil, err.Error
	}

	// record user activity
	actReq := dbmodels.UserActivityRequest{
		UserID:       auther.ID,
		OrgUID:       auther.OrgUID,
		Action:       constants.LoginAction,
		ObjectID:     null.Int64From(auther.ID),
		ObjectType:   null.StringFrom(string(constants.AutherObject)),
		SessionToken: auther.SessionToken,
	}
	_, err = r.services.UserActivityService.Create(ctx, tx, actReq)
	if err != nil {
		return nil, err.Error
	}

	// commit db transaction
	if err := r.services.DBTX.CommitTx(ctx, tx); err != nil {
		return nil, err.Error
	}

	return auther, nil
}

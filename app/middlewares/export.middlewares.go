package middlewares

import (
	"context"

	"github.com/gofrs/uuid"
)

// GetSessionToken reads authToken from context and returns auther by decoding the token
func GetSessionToken(ctx context.Context) *uuid.UUID {
	authToken := authTokenFromContext(ctx)

	token, err := uuid.FromString(authToken)
	if err != nil {
		return nil
	}
	return &token
}

// GetOrgUID reads and returns orgUID from context
func GetOrgUID(ctx context.Context) *uuid.UUID {
	uidStr := orgUIDFromContext(ctx)
	if uidStr == "" {
		return nil
	}
	uid, err := uuid.FromString(uidStr)
	if err != nil {
		return nil
	}
	return &uid
}

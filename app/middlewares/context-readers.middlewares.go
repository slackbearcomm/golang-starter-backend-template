package middlewares

import "context"

// AuthTokenFromContext finds the user from the context. REQUIRES Middleware to have run.
func authTokenFromContext(ctx context.Context) string {
	authToken := ctx.Value(authTokenCtxKey).(string)
	return authToken
}

// orgUIDFromContext finds the user from the context. REQUIRES Middleware to have run.
func orgUIDFromContext(ctx context.Context) string {
	return ctx.Value(orgCtxKey).(string)
}

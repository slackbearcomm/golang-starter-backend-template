package middlewares

type contextKey struct {
	name string
}

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses

var authTokenCtxKey = &contextKey{"auth_token_ctx"}
var orgCtxKey = &contextKey{"org_ctx"}

package handlers

import (
	"gogql/app/services"
	"gogql/app/store/filestore"
)

type Handlers struct {
	AuthHandler    *AuthHandler
	GraphQLHandler *GraphQLHandler
}

func NewHandlers(s *services.Services, fs *filestore.FileStore) *Handlers {
	return &Handlers{
		NewAuthHandler(s),
		NewGraphQLHandler(s, fs),
	}
}

package resolvers

import (
	"context"
	"gogql/app/models"
	"gogql/app/models/dbmodels"
	"gogql/utils/faulterr"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/h2non/filetype"
)

func (r *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*dbmodels.File, error) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UploadFile); err != nil {
		return nil, err.Error
	}

	obj, uploadErr := r.Upload(ctx, file, auther)
	if uploadErr != nil {
		return nil, uploadErr
	}
	return obj, nil
}

func (r *mutationResolver) FileUploadMultiple(ctx context.Context, files []graphql.Upload) ([]dbmodels.File, error) {
	auther, err := r.GetAuther(ctx)
	if err != nil {
		return nil, err.Error
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.UploadFile); err != nil {
		return nil, err.Error
	}

	objects := []dbmodels.File{}
	for _, file := range files {
		obj, uploadErr := r.Upload(ctx, file, auther)
		if uploadErr != nil {
			return nil, uploadErr
		}

		objects = append(objects, *obj)
	}
	return objects, nil
}

// Upload adds a new blob to the db
func (r *mutationResolver) Upload(ctx context.Context, file graphql.Upload, auther *models.Auther) (*dbmodels.File, error) {
	// Read file data
	fileObj, readErr := io.ReadAll(file.File)
	if readErr != nil {
		return nil, faulterr.NewBadRequestError("file upload - read file").Error
	}

	// get mime type
	kind, matchErr := filetype.Match(fileObj)
	if matchErr != nil {
		return nil, faulterr.NewBadRequestError("file upload - get mime type").Error
	}

	if kind == filetype.Unknown {
		return nil, faulterr.NewBadRequestError("file upload - image type is unknown").Error
	}

	// mimeType := kind.MIME.Value
	// extension := kind.Extension

	obj, err := r.filestore.UploadFile(file.Filename, fileObj)
	if err != nil {
		return nil, err.Error
	}
	return obj, nil
}

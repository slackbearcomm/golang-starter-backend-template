package faulterr

// NewBadRequestError structure
func NewBadRequestError(msg string) *FaultErr {
	var err error
	return badRequestErr(msg, err)
}

// NewFrobiddenError structure
func NewFrobiddenError(msg string) *FaultErr {
	var err error
	return frobiddenErr(msg, err)
}

// NewNotAcceptableError structure
func NewNotAcceptableError(msg string) *FaultErr {
	var err error
	return notAcceptableErr(msg, err)
}

// NewUnauthorizedError structure
func NewUnauthorizedError(msg string) *FaultErr {
	var err error
	return unauthorizedErr(msg, err)
}

// NewNotFoundError structure
func NewNotFoundError(msg string) *FaultErr {
	var err error
	return notFoundErr(msg, err)
}

// NewUnprocessableEntityError structure
func NewUnprocessableEntityError(msg string) *FaultErr {
	var err error
	return unprocessableEntityErr(msg, err)
}

// NewInternalServerError structure
func NewInternalServerError(msg string) *FaultErr {
	var err error
	return internalServerErr(msg, err)
}

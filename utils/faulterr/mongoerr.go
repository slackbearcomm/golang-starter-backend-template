package faulterr

// NewMongoError structure
func NewMongoError(err error, msg string) *FaultErr {
	switch err.Error() {
	case "mongo: no documents in result":
		return notFoundErr(msg, err)
	default:
		msg := "Something went wrong, please try again"
		return internalServerErr(msg, err)
	}
}

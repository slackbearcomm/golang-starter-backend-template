package faulterr

// NewPostgresError structure
func NewPostgresError(err error, msg string) *FaultErr {
	switch err.Error() {
	case "no rows in result set":
		return notFoundErr(msg, err)
	default:
		msg := "Something went wrong, please try again"
		return internalServerErr(msg, err)
	}
}

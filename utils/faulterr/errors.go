package faulterr

import (
	"fmt"
	"gogql/utils/logger"
	"net/http"
)

// badRequestErr structure
func badRequestErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Status:  http.StatusBadRequest,
		Error:   fmt.Errorf(msg),
		Message: msg,
	}
}

// frobiddenErr structure
func frobiddenErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Status:  http.StatusForbidden,
		Error:   fmt.Errorf(msg),
		Message: msg,
	}
}

// unauthorizedErr structure
func unauthorizedErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Status:  http.StatusUnauthorized,
		Error:   fmt.Errorf(msg),
		Message: msg,
	}
}

// notFoundErr structure
func notFoundErr(msg string, err error) *FaultErr {
	// logger.Error(err, msg)
	return &FaultErr{
		Status:  http.StatusNotFound,
		Error:   fmt.Errorf("not found: %s", msg),
		Message: msg,
	}
}

// notAcceptableErr structure
func notAcceptableErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Status:  http.StatusNotAcceptable,
		Error:   fmt.Errorf(msg),
		Message: msg,
	}
}

// unprocessableEntityErr structure
func unprocessableEntityErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Status:  http.StatusUnprocessableEntity,
		Error:   fmt.Errorf(msg),
		Message: msg,
	}
}

// internalServerErr structure
func internalServerErr(msg string, err error) *FaultErr {
	logger.Error(err, msg)
	return &FaultErr{
		Status:  http.StatusInternalServerError,
		Error:   fmt.Errorf(msg),
		Message: msg,
	}
}

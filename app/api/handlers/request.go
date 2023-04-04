package handlers

import (
	"gogql/utils/faulterr"
	"strconv"
)

func ConvertStrToInt64(idParam string) (int64, *faulterr.FaultErr) {
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, faulterr.NewBadRequestError("id should be a number")
	}
	return id, nil
}

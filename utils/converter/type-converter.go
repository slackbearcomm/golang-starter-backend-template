package converter

import (
	"gogql/utils/faulterr"
	"strconv"
)

func StrToInt64(numStr string) (int64, *faulterr.FaultErr) {
	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return 0, faulterr.NewBadRequestError("num string should be a number")
	}
	return num, nil
}

func StrToFloat64(numStr string) (float64, *faulterr.FaultErr) {
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, faulterr.NewBadRequestError("num string should be a number")
	}
	return num, nil
}

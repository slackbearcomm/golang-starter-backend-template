package helpers

import (
	"encoding/hex"
	"errors"
	"fmt"
	"gogql/utils/faulterr"
	"regexp"
	"time"
)

// holds functions to be used everywhere

var regexEmailAddr = regexp.MustCompile(`^.+@.+\\..+$`)

// CheckValidEmail basic check to make sure email address is valid, return nil if good
func CheckValidEmail(str string) error {
	if len(str) == 0 {
		return fmt.Errorf("IsValidEmail len: %w", errors.New("email cannot be blank"))
	}

	if !regexEmailAddr.MatchString(str) {
		return fmt.Errorf("IsValidEmail match: %w", errors.New("invalid email address string"))
	}

	return nil
}

var regexProcoreAPIFormat = regexp.MustCompile("[a-f0-9]{64}")

func CheckValidProcoreAPI(str string) error {
	if !regexProcoreAPIFormat.MatchString(str) {
		return fmt.Errorf("invalid Procore api matched: %w", errors.New("invalid procore api format"))
	}
	return nil
}

// StringSliceExist find out if string exist in string slice
func StringSliceExist(arr []string, match string) bool {
	for _, str := range arr {
		if str == match {
			return true
		}
	}

	return false
}

// StringSlicePos find out position of matched string in slice, -1 is not found
func StringSlicePos(arr []string, match string) int {
	for i, str := range arr {
		if str == match {
			return i
		}
	}

	return -1
}

// StringSliceReverse return reverse ordered string slice
func StringSliceReverse(arr []string) []string {
	l := len(arr)
	strs := make([]string, l)

	for i, str := range arr {
		strs[l-i-1] = str
	}

	return strs
}

// StringSliceUnique returns non-repeated string slice
func StringSliceUnique(arr []string) []string {
	newArr := []string{}

	for _, str := range arr {
		if !StringSliceExist(newArr, str) {
			newArr = append(newArr, str)
		}
	}

	return newArr
}

// StripSensitive remove sensitive data from string
func StripSensitive(txt string) string {
	txt2 := StripPassword(txt)

	return txt2
}

// StripPassword require these regexs to work
var (
	reSP1 = regexp.MustCompile(`(?i)("password"):".+?"`)
	reSP2 = regexp.MustCompile(`(?i)("password_confirm"):".+?"`)
	reSP3 = regexp.MustCompile(`(?i)("passwd"):".+?"`)
	reSP4 = regexp.MustCompile(`(?i)("pass"):".+?"`)
	reSP5 = regexp.MustCompile(`(?i)("pwd"):".+?"`)
	reSP6 = regexp.MustCompile(`(?i)("pw"):".+?"`)
	reSP7 = regexp.MustCompile(`password`)
)

// StripPassword remove password from string
func StripPassword(txt string) string {
	bRepl := []byte(`$1:"**REMOVED**"`)

	// duplicate
	str := txt

	// hide password
	str = string(reSP1.ReplaceAll([]byte(str), bRepl))

	// hide password
	str = string(reSP2.ReplaceAll([]byte(str), bRepl))

	// other matches, just incase if someone uses non-standard name

	str = string(reSP3.ReplaceAll([]byte(str), bRepl))

	str = string(reSP4.ReplaceAll([]byte(str), bRepl))

	str = string(reSP5.ReplaceAll([]byte(str), bRepl))

	str = string(reSP6.ReplaceAll([]byte(str), bRepl))

	// rename password, just in case the sentry decides to filter it
	str = string(reSP7.ReplaceAll([]byte(str), []byte("paxxword")))

	return str
}

// GetCartonSpreadsheetLink generates carton spread sheet link/url
func GetCartonSpreadsheetLink(code string, cartons int, cartonCount int) string {
	startCode := ""
	endCode := ""

	for i := 0; i < cartons; i++ {
		code := fmt.Sprintf("CAR%05d", cartonCount)

		if i == 0 {
			startCode = code
		} else if i == cartons-1 {
			endCode = code
		}

		cartonCount++
	}

	return fmt.Sprintf("/api/files/sheet?type=carton&from=%s&to=%s", startCode, endCode)
}

// HexToBytes convert hex strings to bytes
func HexToBytes(strHex string) ([]byte, error) {
	decoded, err := hex.DecodeString(strHex)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func ValidateTokenExpiry(expiresAt time.Time) *faulterr.FaultErr {
	if expiresAt.UTC().Unix() < time.Now().UTC().Unix() {
		return faulterr.NewUnauthorizedError("session is expired")
	}
	return nil
}

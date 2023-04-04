package models

import (
	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

type Auther struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	IsAdmin      bool          `json:"IsAdmin"`
	OrgUID       uuid.NullUUID `json:"orgaUID"`
	RoleID       null.Int64    `json:"roleID"`
	SessionToken uuid.UUID     `json:"sessionToken"`
}

type OTPRequest struct {
	Email null.String `json:"email"`
	Phone null.String `json:"phone"`
}

type LoginRequest struct {
	Email null.String `json:"email"`
	Phone null.String `json:"phone"`
	OTP   string      `json:"otp"`
}

// ValueToken struct
type ValueToken struct {
	TokenString string `json:"tokenString"`
}

// TokenHeader struct
type TokenHeader struct {
	TYP string `json:"typ"`
	ALG string `json:"alg"`
}

// TokenPayload struct
type TokenPayload struct {
	Auther
	Authorized bool  `json:"authorized"`
	Expiry     int64 `json:"exp"`
}

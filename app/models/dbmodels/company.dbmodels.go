package dbmodels

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/volatiletech/null"
)

///////////////////
//   DB Models   //
///////////////////

type Permission struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Organization struct {
	ID         int64       `json:"id"`
	UID        uuid.UUID   `json:"uid"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Website    null.String `json:"website"`
	Logo       File        `json:"logo"`
	Sector     string      `json:"sector"`
	Status     string      `json:"status"`
	IsFinal    bool        `json:"isFinal"`
	IsArchived bool        `json:"isArchived"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

type Department struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	OrgUID     uuid.UUID `json:"orgUID"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	IsFinal    bool      `json:"isFinal"`
	IsArchived bool      `json:"isArchived"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Role struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	OrgUID       uuid.UUID `json:"orgUID"`
	DepartmentID int64     `json:"departmentID"`
	Name         string    `json:"name"`
	Permissions  []string  `json:"permissions"`
	IsManagement bool      `json:"isManagement"`
	Status       string    `json:"status"`
	IsFinal      bool      `json:"isFinal"`
	IsArchived   bool      `json:"isArchived"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type User struct {
	ID         int64         `json:"id"`
	FirstName  string        `json:"firstName"`
	LastName   string        `json:"lastName"`
	Email      string        `json:"email"`
	Phone      string        `json:"phone"`
	IsAdmin    bool          `json:"isAdmin"`
	OrgUID     uuid.NullUUID `json:"orgUID"`
	RoleID     null.Int64    `json:"roleID"`
	Status     string        `json:"status"`
	IsFinal    bool          `json:"isFinal"`
	IsArchived bool          `json:"isArchived"`
	CreatedAt  time.Time     `json:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt"`
}

type OTPSession struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userID"`
	Token     string    `json:"token"`
	IsValid   bool      `json:"isValid"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type AuthSession struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userID"`
	Token     uuid.UUID `json:"token"`
	IsValid   bool      `json:"isValid"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserActivity struct {
	ID           int64         `json:"id"`
	UserID       int64         `json:"userID"`
	OrgUID       uuid.NullUUID `json:"orgUID"`
	Action       string        `json:"action"`
	ObjectID     null.Int64    `json:"objectID"`
	ObjectType   null.String   `json:"objectType"`
	SessionToken uuid.UUID     `json:"sessiosToken"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

////////////////////////
//   Request Models   //
////////////////////////

type OrganizationRequest struct {
	Name       string      `json:"name"`
	Website    null.String `json:"website"`
	Logo       File        `json:"logo"`
	Sector     string      `json:"sector"`
	Status     string      `json:"status"`
	IsArchived bool        `json:"isArchived"`
}

type OrganizationRegisterRequest struct {
	OrgName   string      `json:"orgName"`
	Website   null.String `json:"website"`
	Logo      File        `json:"logo"`
	Sector    string      `json:"sector"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Email     string      `json:"email"`
	Phone     string      `json:"phone"`
}

type DepartmentRequest struct {
	OrgUID  uuid.UUID `json:"orgUID"`
	Name    string    `json:"name"`
	IsFinal bool      `json:"isFinal"`
}

type RoleRequest struct {
	OrgUID       uuid.UUID `json:"orgUID"`
	DepartmentID int64     `json:"departmentID"`
	Name         string    `json:"name"`
	Permissions  []string  `json:"permissions"`
	IsManagement bool      `json:"isManagement"`
	Status       string    `json:"status"`
	IsFinal      bool      `json:"isFinal"`
}

type SuperAdminRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
}

type UserRequest struct {
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Email     string        `json:"email"`
	Phone     string        `json:"phone"`
	IsAdmin   bool          `json:"isAdmin"`
	OrgUID    uuid.NullUUID `json:"orgUID"`
	RoleID    null.Int64    `json:"roleID"`
	Status    string        `json:"status"`
	IsFinal   bool          `json:"isFInal"`
}

type UserActivityRequest struct {
	UserID       int64         `json:"userID"`
	OrgUID       uuid.NullUUID `json:"orgUID"`
	Action       string        `json:"action"`
	ObjectID     null.Int64    `json:"objectID"`
	ObjectType   null.String   `json:"objectType"`
	SessionToken uuid.UUID     `json:"sessiosToken"`
}

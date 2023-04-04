package authtoken

import (
	"gogql/app/models"
	"gogql/utils/faulterr"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type Payload struct {
	TokenString string
	ExpiresAt   time.Time
}

type Claims struct {
	models.Auther
	jwt.StandardClaims
}

func Generate(auther *models.Auther) (*Payload, *faulterr.FaultErr) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(24) * time.Duration(30))
	claims := &Claims{
		Auther: *auther,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, faulterr.NewBadRequestError("bad token request")
	}

	payload := &Payload{
		TokenString: tokenString,
		ExpiresAt:   expirationTime,
	}

	return payload, nil
}

func Decode(tokenString string) (*models.Auther, *faulterr.FaultErr) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, faulterr.NewUnauthorizedError("unauthorized request: invalid signature")
		}
		return nil, faulterr.NewBadRequestError(err.Error())
	}
	if !token.Valid {
		return nil, faulterr.NewUnauthorizedError("unauthorized request: token is not valid")
	}

	auther := &models.Auther{
		ID:      claims.ID,
		Name:    claims.Name,
		IsAdmin: claims.IsAdmin,
		OrgUID:  claims.OrgUID,
		RoleID:  claims.RoleID,
	}

	return auther, nil
}

func keyFunc(t *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}

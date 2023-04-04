package handlers

import (
	"context"
	"gogql/app/models"
	"gogql/app/services"
	"gogql/utils/authtoken"
	"gogql/utils/faulterr"
	"net/http"
	"time"
)

type AuthHandler struct {
	services *services.Services
	// jwt      *authentication.JWT
}

func NewAuthHandler(s *services.Services) *AuthHandler {
	// jwt := authentication.NewJWT(jwt.SecretKey, 30)
	// return &AuthHandler{s, jwt}
	return &AuthHandler{s}
}

// GenerateToken generates a jwt token and sets cookie
func (h *AuthHandler) GenerateToken(w http.ResponseWriter, ctx context.Context, auther *models.Auther) (*http.Cookie, *faulterr.FaultErr) {
	tokenPayload, err := authtoken.Generate(auther)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    tokenPayload.TokenString,
		Expires:  tokenPayload.ExpiresAt,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	return cookie, nil
}

// ListPermissions Handler
func (h *AuthHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	result := models.ListPermissions()

	response := ResponseBody{
		Data:       result,
		Message:    "Permissions list",
		StatusCode: http.StatusOK,
	}

	RestResponse(w, r, response.StatusCode, response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w,
		&http.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})

	response := ResponseBody{
		Data:       nil,
		Message:    "user logout",
		StatusCode: http.StatusCreated,
	}

	RestResponse(w, r, response.StatusCode, response)
}

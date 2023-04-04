package handlers

import (
	"gogql/app/models"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

// ResponseBody struct
type ResponseBody struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"status"`
}

type AuthData struct {
	CookieToken *http.Cookie   `json:"cookieToken"`
	Auther      *models.Auther `json:"auther"`
	Token       string         `json:"token"`
}

// RestResponse handles the http status and renders body in JSON
func RestResponse(w http.ResponseWriter, r *http.Request, status int, body interface{}) {
	render.Status(r, status)
	render.JSON(w, r, body)
}

func SetCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   false,
	}
	http.SetCookie(w, &cookie)
}

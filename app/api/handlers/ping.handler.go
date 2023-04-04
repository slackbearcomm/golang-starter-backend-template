package handlers

import (
	"net/http"
)

// Ping test
func Ping(w http.ResponseWriter, r *http.Request) {
	response := ResponseBody{
		Data:       "pong",
		Message:    "Ping Response",
		StatusCode: http.StatusOK,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// King test
func King(w http.ResponseWriter, r *http.Request) {
	response := ResponseBody{
		Data:       "kong",
		Message:    "King Response",
		StatusCode: http.StatusOK,
	}

	RestResponse(w, r, response.StatusCode, response)
}

// Ding test
func Ding(w http.ResponseWriter, r *http.Request) {
	response := ResponseBody{
		Data:       "dong",
		Message:    "Ding Response",
		StatusCode: http.StatusOK,
	}

	RestResponse(w, r, response.StatusCode, response)
}

package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version,omitempty"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "UP",
	}
	
	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
package server

import (
	"net/http"
	"ozon_service/service"

	"github.com/gorilla/mux"
)

func New(s *service.Service) http.Handler {
	h := Handler{svc: s}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/job", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/job", h.GetByURL).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/job", h.DeleteByURL).Methods(http.MethodDelete)
	return r
}

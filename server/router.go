package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	server http.Server
}

func New() *Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/create", Create).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/get_by_url", GetByURL).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/delete_by_url", DeleteByURL).Methods(http.MethodDelete)
	return &Router{
		server: http.Server{
			Addr:    ":8080",
			Handler: r,
		}}
}

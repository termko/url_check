package server

import (
	"net/http"
	"ozon_service/service"
)

type Handler struct {
	svc *service.Service
}

func Create(w http.ResponseWriter, r *http.Request) {

}

func GetByURL(w http.ResponseWriter, r *http.Request) {

}

func DeleteByURL(w http.ResponseWriter, r *http.Request) {

}

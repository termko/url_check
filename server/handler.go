package server

import (
	"encoding/json"
	"net/http"
	"ozon_service/service"
	"strconv"
	"time"
)

type Handler struct {
	svc *service.Service
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var job service.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong input format"))
		return
	}
	if job.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong input format"))
		return
	}
	if job.Interval == 0 {
		job.Interval = service.Duration(24 * time.Hour)
	}
	err = h.svc.Create(r.Context(), &job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetByURL(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	URL, ok := query["url"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must provide 1 URL"))
		return
	}
	if len(URL) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must provide 1 URL"))
		return
	}
	result, err := h.svc.GetByURL(r.Context(), URL[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	data, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *Handler) DeleteByURL(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	URL, ok := query["url"]
	if !ok || len(URL) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must provide 1 URL"))
		return
	}
	err := h.svc.DeleteByURL(r.Context(), URL[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetScore(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	dateStrings, ok := query["date"]
	if !ok || len(dateStrings) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must provide a date parameter"))
		return
	}
	nStrings, ok := query["n"]
	if !ok || len(nStrings) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must provide a score parameter (n)"))
		return
	}
	date, err := time.Parse("2006-01-02", dateStrings[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong format of date! (YYYY-MM-DD)"))
		return
	}
	n, err := strconv.Atoi(nStrings[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parameter n must be number!"))
		return
	}
	score := &service.Score{Date: date, N: n}
	result, err := h.svc.GetScore(r.Context(), score)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	data, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

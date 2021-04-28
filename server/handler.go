package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ozon_service/service"
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
	if job.Interval == "" || job.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Wrong input format"))
		return
	}
	err = h.svc.Create(r.Context(), &job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something crashed :D"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) GetByURL(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	URL, ok := query["url"]
	if ok {
		if len(URL) == 1 {
			result, err := h.svc.GetByURL(r.Context(), URL[0])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Something broke ;)"))
				return
			}
			for _, i := range result {
				fmt.Println(i.PingTime, i.StatusCode)
			}
			if len(result) == 0 {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("You must provide 1 URL"))
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must provide 1 URL"))
		return
	}
}

func (h *Handler) DeleteByURL(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	URL, ok := query["url"]
	if ok {
		if len(URL) == 1 {
			err := h.svc.DeleteByURL(r.Context(), URL[0])
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Something doesn't work ;("))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("You must provide 1 URL"))
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You must provide 1 URL"))
		return
	}
}

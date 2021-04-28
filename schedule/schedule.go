package schedule

import "ozon_service/service"

type Schedule struct {
	job *service.Job
}

func New(j *service.Job) *Schedule {
	return &Schedule{job: j}
}

func PingURL(job *service.Job) error {
	return nil
}

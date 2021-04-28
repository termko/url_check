package service

import (
	"context"
	"fmt"
	"time"
	
	"github.com/robfig/cron/v3"
)

type Job struct {
	URL      string `json:"url"`
	Interval string `json:"interval"`
}

type JobResult struct {
	StatusCode int
	PingTime   time.Time
}

type Service struct {
	repo      Repository
	scheduler *cron.Cron
}

type Repository interface {
	Create(job *Job, jobID cron.EntryID) error
	GetByURL(URL string) ([]JobResult, error)
	DeleteByURL(URL string) error
	GetIDByURL(URL string) (cron.EntryID, error)
}

func New(repo Repository, scheduler *cron.Cron) *Service {
	return &Service{repo: repo, scheduler: scheduler}
}

func (s *Service) Create(ctx context.Context, job *Job) error {
	jobID, err := s.scheduler.AddFunc("@every "+job.Interval, func() { fmt.Println(job.Interval, job.URL) })

	if err != nil {
		return err
	}
	err = s.repo.Create(job, jobID)
	if err != nil {
		s.scheduler.Remove(jobID)
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) GetByURL(ctx context.Context, URL string) ([]JobResult, error) {
	result, err := s.repo.GetByURL(URL)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) DeleteByURL(ctx context.Context, URL string) error {
	jobID, err := s.repo.GetIDByURL(URL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.scheduler.Remove(jobID)
	s.repo.DeleteByURL(URL)
	return nil
}

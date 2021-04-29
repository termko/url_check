package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/robfig/cron/v3"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type Job struct {
	URL      string   `json:"url"`
	Interval Duration `json:"interval"`
}

type Score struct {
	Date time.Time
	N    int
}

type ScoreResult struct {
	URL   string
	Score int
}

type JobResult struct {
	StatusCode int
	PingTime   time.Time
	Available  bool
}

type Service struct {
	repo       Repository
	scheduler  *cron.Cron
	jobIDToURL map[string]cron.EntryID
}

type Repository interface {
	Create(ctx context.Context, job *Job, jobID cron.EntryID) error
	GetByURL(ctx context.Context, URL string) ([]JobResult, error)
	DeleteByURL(ctx context.Context, URL string) error
	CheckURL(ctx context.Context, URL string) error
	InsertPing(ctx context.Context, URL string, jobResult JobResult) error
	GetScore(ctx context.Context, score *Score) ([]ScoreResult, error)
}

func New(repo Repository, scheduler *cron.Cron) *Service {
	return &Service{repo: repo, scheduler: scheduler, jobIDToURL: make(map[string]cron.EntryID)}
}

func (s *Service) Create(ctx context.Context, job *Job) error {
	jobID, err := s.scheduler.AddFunc("@every "+time.Duration(job.Interval).String(), func() {
		ctx := context.Background()
		err := s.repo.CheckURL(ctx, job.URL)
		if err != nil {
			jobID, ok := s.jobIDToURL[job.URL]
			if ok {
				delete(s.jobIDToURL, job.URL)
				s.scheduler.Remove(jobID)
				return
			}
		}
		pingTime := time.Now()
		status := 0
		available := true
		client := http.Client{
			Timeout: 5 * time.Second,
		}
		resp, err := client.Get(job.URL)
		if err != nil || resp.StatusCode >= 400 {
			available = false
		} else {
			status = resp.StatusCode
		}
		jobResult := JobResult{StatusCode: status, PingTime: pingTime, Available: available}
		err = s.repo.InsertPing(ctx, job.URL, jobResult)
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		return err
	}
	err = s.repo.Create(ctx, job, jobID)
	if err != nil {
		s.scheduler.Remove(jobID)
		fmt.Println(err)
		return err
	}
	s.jobIDToURL[job.URL] = jobID
	return nil
}

func (s *Service) GetByURL(ctx context.Context, URL string) ([]JobResult, error) {
	result, err := s.repo.GetByURL(ctx, URL)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Service) DeleteByURL(ctx context.Context, URL string) error {
	if jobID, ok := s.jobIDToURL[URL]; ok {
		s.scheduler.Remove(jobID)
	}
	s.repo.DeleteByURL(ctx, URL)
	return nil

}

func (s *Service) GetScore(ctx context.Context, score *Score) ([]ScoreResult, error) {
	result, err := s.repo.GetScore(ctx, score)
	if err != nil {
		return nil, err
	}
	return result, nil
}

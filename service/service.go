package service

import "time"

type URL struct {
	Url      string        `json:"url"`
	Interval time.Duration `json:"interval"`
}

type Service struct {
	repo Repository
}

type Repository interface {
	Create(url string, interval time.Duration) error
	GetByURL(url string) (*URL, error)
	DeleteByURL(url string) error
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(urlstr string, interval time.Duration) (*URL, error) {
	err := s.repo.Create(urlstr, interval)
	if err != nil {
		return nil, err
	}
	url, err := s.repo.GetByURL(urlstr)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (s *Service) GetByURL(url string) (*URL, error) {
	return s.repo.GetByURL(url)
}

func (s *Service) DeleteByURL(urlstr string) (*URL, error) {
	url, err := s.repo.GetByURL(urlstr)
	if err != nil {
		return nil, err
	}
	err = s.repo.DeleteByURL(urlstr)
	if err != nil {
		return nil, err
	}
	return url, nil
}

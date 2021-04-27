package repository

import (
	"database/sql"
	"ozon_service/service"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func New(DB *sql.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(urlstr string, interval time.Duration) (*service.URL, error) {
	query := `INSERT INTO checks (URL) VALUES (?)`
	_, err := r.DB.Exec(query, urlstr)
	if err != nil {
		return nil, err
	}
	url, err := r.GetByURL(urlstr)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (r *Repository) GetByURL(urlstr string) (*service.URL, error) {
	query := `SELECT URL, Interval FROM checks WHERE URL=?`
	row := r.DB.QueryRow(query, urlstr)
	var url *service.URL
	err := row.Scan(url.Url, url.Interval)
	if err != nil {
		return nil, err
	}
	return url, nil
}

func (r *Repository) DeleteByURL(urlstr service.URL) error {
	query := `DELETE FROM checks WHERE URL=?`
	_, err := r.DB.Exec(query, urlstr)
	if err != nil {
		return err
	}
	return nil
}

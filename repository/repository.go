package repository

import (
	"database/sql"
	"ozon_service/service"

	"github.com/robfig/cron/v3"
)

type Repository struct {
	DB *sql.DB
}

func New(DB *sql.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(job *service.Job, jobID cron.EntryID) error {
	query := `INSERT INTO schedule (URL, jobID) VALUES ($1, $2)`
	_, err := r.DB.Exec(query, job.URL, jobID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByURL(URL string) ([]service.JobResult, error) {
	query := `SELECT statuscode, pingtime FROM ping WHERE URL=$1`
	rows, err := r.DB.Query(query, URL)
	if err != nil {
		return nil, err
	}
	result := make([]service.JobResult, 0)
	for rows.Next() {
		var url service.JobResult
		err := rows.Scan(&(url.StatusCode), &(url.PingTime))
		if err != nil {
			return nil, err
		}
		result = append(result, url)
	}

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) DeleteByURL(URL string) error {
	query := `DELETE FROM schedule WHERE URL=$1`
	_, err := r.DB.Exec(query, URL)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetIDByURL(URL string) (cron.EntryID, error) {
	query := `SELECT jobID FROM schedule where URL=$1`
	var jobID cron.EntryID
	err := r.DB.QueryRow(query, URL).Scan(&jobID)
	if err != nil {
		return -1, err
	}
	return jobID, nil
}

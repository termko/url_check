package repository

import (
	"context"
	"database/sql"
	"fmt"
	"ozon_service/service"

	"github.com/robfig/cron/v3"
)

type Repository struct {
	DB *sql.DB
}

func New(DB *sql.DB) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) Create(ctx context.Context, job *service.Job, jobID cron.EntryID) error {
	query := `INSERT INTO schedule (URL, jobID) VALUES ($1, $2)`
	_, err := r.DB.ExecContext(ctx, query, job.URL, jobID)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByURL(ctx context.Context, URL string) ([]service.JobResult, error) {
	query := `SELECT statuscode, pingtime, available FROM ping WHERE URL=$1`
	rows, err := r.DB.QueryContext(ctx, query, URL)
	if err != nil {
		return nil, err
	}
	var result []service.JobResult
	for rows.Next() {
		var url service.JobResult
		err := rows.Scan(&url.StatusCode, &url.PingTime, &url.Available)
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

func (r *Repository) DeleteByURL(ctx context.Context, URL string) error {
	query := `DELETE FROM schedule WHERE url=$1`
	_, err := r.DB.ExecContext(ctx, query, URL)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CheckURL(ctx context.Context, URL string) error {
	query := `SELECT url FROM schedule WHERE url=$1`
	var check string
	err := r.DB.QueryRowContext(ctx, query, URL).Scan(&check)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) InsertPing(ctx context.Context, URL string, jobResult service.JobResult) error {
	query := `INSERT INTO ping (url, statuscode, pingtime, available) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.ExecContext(ctx, query, URL, jobResult.StatusCode, jobResult.PingTime, jobResult.Available)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetScore(ctx context.Context, score *service.Score) ([]service.ScoreResult, error) {
	query := `SELECT url, count(available) FROM ping WHERE pingtime >= $1 AND available = true GROUP BY url`
	rows, err := r.DB.QueryContext(ctx, query, score.Date)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var scoreResult []service.ScoreResult
	for rows.Next() {
		var result service.ScoreResult
		err = rows.Scan(&result.URL, &result.Score)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if result.Score >= score.N {
			scoreResult = append(scoreResult, result)
		}
	}
	return scoreResult, nil
}

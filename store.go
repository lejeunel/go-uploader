package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var job_schema = `
CREATE TABLE IF NOT EXISTS jobs (
    id varchar(16),
    uri_source text,
    uri_destination text,
    status tinyint,
	created_at text,
	updated_at text
);`

var transaction_schema = `
CREATE TABLE IF NOT EXISTS transactions (
    id varchar(16),
    job_id varchar(16),
    uri_source text,
    uri_destination text,
	status tinyint,
	created_at text,
	updated_at text
);`

type Store interface {
	AppendJob(job *Job) (*Job, error)
	AppendJobTransactions(job *Job) (*Job, error)
	UpdateJobStatus(job *Job) error
	UpdateTransactionStatus(t *Transaction) error
	FindJob(uriSource string, uriDestination string) (*Job, error)
	getTransactions(jobID uuid.UUID) ([]Transaction, error)
	DeleteJob(job *Job) error
}

type SQLiteStore struct {
	store Store
	db    *sqlx.DB
}

func (s *SQLiteStore) AppendJob(job *Job) (*Job, error) {
	job.ID = uuid.New()
	query := "INSERT INTO jobs (id, uri_source, uri_destination, status, created_at) VALUES (?, ?, ?, ?, ?)"
	_, err := s.db.Exec(query, job.ID, job.UriSource, job.UriDestination, job.Status, time.Now().String())

	return job, err
}

func (s *SQLiteStore) AppendJobTransactions(job *Job) (*Job, error) {
	var transactions []Transaction
	tx := s.db.MustBegin()
	for _, t := range job.Transactions {
		t.ID = uuid.New()
		t.JobId = job.ID
		t.CreatedAt = time.Now().String()
		query := "INSERT INTO transactions (id, job_id, uri_source, uri_destination, status, created_at) VALUES (?, ?, ?, ?, ?, ?)"
		_, err := tx.Exec(query, t.ID, job.ID, t.UriSource, t.UriDestination, t.Status, t.CreatedAt)
		if err != nil {
			return job, err
		}
		transactions = append(transactions, t)
	}
	tx.Commit()
	job.Transactions = transactions
	return job, nil
}

func (s *SQLiteStore) UpdateJobStatus(job *Job) error {
	query := "UPDATE jobs SET status=?, updated_at=? WHERE id=?"
	_, err := s.db.Exec(query, job.Status, time.Now().String(), job.ID)

	return err
}

func (s *SQLiteStore) UpdateTransactionStatus(t *Transaction) error {

	query := "UPDATE transactions SET status=?, updated_at=? WHERE id=?"
	_, err := s.db.Exec(query, t.Status, time.Now().String(), t.ID)
	return err
}

func (s *SQLiteStore) getTransactions(jobID uuid.UUID) ([]Transaction, error) {
	query := "SELECT id, job_id, status, uri_source, uri_destination, created_at, updated_at FROM transactions WHERE job_id=?"
	transactions := []Transaction{}
	err := s.db.Select(&transactions, query, jobID)

	return transactions, err
}

func (s *SQLiteStore) FindJob(uriSource string, uriDestination string) (*Job, error) {
	query := "SELECT id, uri_source, uri_destination, status, created_at, updated_at FROM jobs WHERE uri_source=? AND uri_destination=?"
	job := Job{}
	err := s.db.Get(&job, query, uriSource, uriDestination)

	if err != nil {
		return nil, &jobNotFoundError{fmt.Sprintf("Could not find job with source and destination: %s %s",
			uriSource, uriDestination)}
	}

	job.Transactions, err = s.getTransactions(job.ID)

	return &job, err
}
func (s *SQLiteStore) DeleteJob(job *Job) error {
	_, err_j := s.db.Exec("DELETE FROM jobs WHERE id=?", job.ID)
	_, err_t := s.db.Exec("DELETE FROM transactions WHERE job_id=?", job.ID)

	return errors.Join(err_j, err_t)
}

func NewStore(db *sqlx.DB) *SQLiteStore {

	db.MustExec(transaction_schema)
	db.MustExec(job_schema)

	return &SQLiteStore{db: db}

}

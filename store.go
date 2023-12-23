package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

var job_schema = `
CREATE TABLE jobs (
    id varchar(16),
    uri_source text,
    uri_destination text,
    status tinyint,
	created_at text,
	updated_at text
);`

var transaction_schema = `
CREATE TABLE transactions (
    id varchar(16),
    job_id varchar(16),
    uri_source text,
    uri_destination text,
	status tinyint,
	created_at text,
	updated_at text
);`

type store interface {
	AppendJob(job *Job) (*Job, error)
	AppendJobTransactions(job *Job) (*Job, error)
	UpdateJob(job *Job) error
	UpdateTransaction(t *Transaction) error
	FindJob(uriSource string, uriDestination string) (*Job, error)
	getTransactions(jobID uuid.UUID) ([]Transaction, error)
}

type SQLiteStore struct {
	store store
	db    *sqlx.DB
}

func (s *SQLiteStore) AppendJob(job *Job) (*Job, error) {
	tx := s.db.MustBegin()
	job.ID = uuid.New()
	query := fmt.Sprintf(
		"INSERT INTO jobs (id, uri_source, uri_destination, status, created_at) VALUES (\"%s\", \"%s\", \"%s\", %s, datetime())",
		job.ID, job.UriSource, job.UriDestination,
		strconv.Itoa(job.Status))
	_, err := tx.Exec(query)

	tx.Commit()
	return job, err
}

func (s *SQLiteStore) AppendJobTransactions(job *Job) (*Job, error) {
	var transactions []Transaction
	tx := s.db.MustBegin()
	for _, t := range job.Transactions {
		t.ID = uuid.New()
		t.JobId = job.ID
		query := fmt.Sprintf(
			"INSERT INTO transactions (id, job_id, uri_source, uri_destination, status, created_at) VALUES (\"%s\", \"%s\", \"%s\", \"%s\", %s, datetime())",
			t.ID, job.ID, t.UriSource, t.UriDestination,
			strconv.Itoa(t.Status))
		_, err := tx.Exec(query)
		if err != nil {
			return job, err
		}
		transactions = append(transactions, t)
	}
	tx.Commit()
	job.Transactions = transactions
	return job, nil
}

func (s *SQLiteStore) UpdateJob(job *Job) error {
	tx := s.db.MustBegin()
	_, err := tx.Exec(
		"UPDATE jobs SET status=$1, updated_at=datetime() WHERE id=$2",
		job.Status, job.ID)
	tx.Commit()

	return err
}

func (s *SQLiteStore) UpdateTransaction(t *Transaction) error {

	tx := s.db.MustBegin()
	_, err := tx.Exec(
		"UPDATE transactions SET status=$1, updated_at=datetime() WHERE id=$2",
		t.Status, t.ID)
	tx.Commit()
	return err
}

func (s *SQLiteStore) getTransactions(jobID uuid.UUID) ([]Transaction, error) {
	query := fmt.Sprintf("SELECT id, job_id, status, uri_source, uri_destination, created_at, updated_at FROM transactions WHERE job_id=\"%s\";", jobID)
	transactions := []Transaction{}
	err := s.db.Select(&transactions, query)

	return transactions, err
}

func (s *SQLiteStore) FindJob(uriSource string, uriDestination string) (*Job, error) {
	query := fmt.Sprintf("SELECT id, uri_source, uri_destination, status, created_at, updated_at FROM jobs WHERE uri_source=\"%s\" AND uri_destination=\"%s\"",
		uriSource, uriDestination)
	jobs := []Job{}
	err := s.db.Select(&jobs, query)

	if len(jobs) == 0 {
		return nil, &jobNotFoundError{fmt.Sprintf("Could not find job with source and destination: %s %s",
			uriSource, uriDestination)}
	}

	job := jobs[0]
	job.Transactions, err = s.getTransactions(job.ID)

	return &job, err
}

func NewStore(db *sqlx.DB) *SQLiteStore {

	// Migrate the schema
	db.MustExec(transaction_schema)
	db.MustExec(job_schema)

	return &SQLiteStore{db: db}

}

func NewMockStore() *SQLiteStore {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(1)
	return NewStore(db)
}

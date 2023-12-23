package main

import (
	"github.com/google/uuid"
)

const (
	created = iota
	parsed  = iota
	done    = iota
)

const (
	pending     = iota
	transferred = iota
)

type Job struct {
	ID             uuid.UUID `db:"id"`
	UriSource      string    `db:"uri_source"`
	UriDestination string    `db:"uri_destination"`
	Status         int       `db:"status"`
	CreatedAt      string    `db:"created_at"`
	UpdatedAt      string    `db:"updated_at"`
	Transactions   []Transaction
}

type Transaction struct {
	ID             uuid.UUID `db:"id"`
	JobId          uuid.UUID `db:"job_id"`
	UriSource      string    `db:"uri_source"`
	UriDestination string    `db:"uri_destination"`
	CreatedAt      string    `db:"created_at"`
	UpdatedAt      string    `db:"updated_at"`
	Status         int       `db:"status"`
}

func (j *Job) numDoneTransactions() int {
	numDone := 0
	for _, t := range j.Transactions {
		if t.Status == transferred {
			numDone++
		}
	}
	return numDone
}

func (j *Job) numTransactions() int {
	return len(j.Transactions)
}

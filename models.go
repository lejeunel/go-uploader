package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	gorm.Model
	ID             uuid.UUID
	UriSource      string
	UriDestination string
	Status         int
	Transactions   []Transaction
}

type Transaction struct {
	gorm.Model
	ID     uuid.UUID
	JobId  uuid.UUID
	uriIn  string
	uriOut string
	Status int
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

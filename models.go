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
	id     uuid.UUID
	JobId  uuid.UUID
	uriIn  string
	uriOut string
	status int
}

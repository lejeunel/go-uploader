package main

import (
	"time"
)

const (
	created = iota
	parsed  = iota
	done    = iota
)

type job struct {
	uriSource      string
	uriDestination string
	createdAt      time.Time
	updatedAt      time.Time
	status         int
	transactions   []*transaction
}

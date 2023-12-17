package main

const (
	initialized = iota
	parsed      = iota
	done        = iota
)

type job struct {
	last_status  int
	transactions []transaction
}

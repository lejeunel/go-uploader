package main

const (
	pending     = iota
	transferred = iota
)

type transaction struct {
	uriIn  string
	uriOut string
	status int
}

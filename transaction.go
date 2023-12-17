package main

const (
	pending     = iota
	transferred = iota
)

type transaction struct {
	uri_in  string
	uri_out string
	status  int
}

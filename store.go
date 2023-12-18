package main

type store interface {
	commit(job *job) bool
}

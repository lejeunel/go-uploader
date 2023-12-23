package main

import (
	"fmt"
)

type jobNotFoundError struct {
	msg string
}

func (e *jobNotFoundError) Error() string {
	return fmt.Sprintf("Error fetching job: %s", e.msg)
}

type duplicateJobError struct {
	job *Job
}

func (e *duplicateJobError) Error() string {
	return fmt.Sprintf(
		"Attempting to create a job, but found duplicate with id %s",
		e.job.ID)
}

type schemeError struct {
	allowed_scheme string
	provided_uri   string
}

func (e *schemeError) Error() string {
	return fmt.Sprintf("Bad scheme. Provided uri %s, allowed %s", e.provided_uri, e.allowed_scheme)
}

type sourceError struct {
	provided_uri string
}

func (e *sourceError) Error() string {
	return fmt.Sprintf("Source URI %s not found/valid.", e.provided_uri)
}

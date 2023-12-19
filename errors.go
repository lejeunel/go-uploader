package main

import (
	"fmt"
)

type jobNotFoundError struct {
}

func (e *jobNotFoundError) Error() string {
	return "Job not found error"
}

type jobInitError struct {
}

func (e *jobInitError) Error() string {
	return "Error initializing job"
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

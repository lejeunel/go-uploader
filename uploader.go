package main

import (
	"fmt"
)

type schemeError struct {
	allowed_scheme string
	provided_uri   string
}

func (e *schemeError) Error() string {
	return fmt.Sprintf("Bad scheme. Provided uri %s, allowed %s", e.provided_uri, e.allowed_scheme)
}

type reader interface {
	read(uri string) []byte
	scan(uri string) ([]string, error)
	check_scheme(uri string) error
}

type writer interface {
	write(bytes []byte, uri string) bool
	check_scheme(uri string) error
}

type uploader struct {
	reader
	writer
}

func (u uploader) upload(input string, output string) bool {
	bytes := u.read(input)
	u.write(bytes, output)
	return true
}

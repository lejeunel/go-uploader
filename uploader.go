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

type sourceError struct {
	provided_uri string
}

func (e *sourceError) Error() string {
	return fmt.Sprintf("Source URI %s not found/valid.", e.provided_uri)
}

type reader interface {
	read(uri string) []byte
	scan(uri string) []string
	checkScheme(uri string) error
	checkExists(uri string) error
}

type writer interface {
	write(bytes []byte, uri string) bool
	checkScheme(uri string) error
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

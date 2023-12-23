package main

import "fmt"

type Reader interface {
	read(uri string) ([]byte, error)
	scan(uri string) ([]string, error)
	checkScheme(uri string) error
	checkExists(uri string) error
}

type readerError struct {
	uri string
}

func (e *readerError) Error() string {
	return fmt.Sprintf("Error reading %s", e.uri)
}

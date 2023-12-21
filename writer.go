package main

import "fmt"

type writer interface {
	write(bytes []byte, uri string) error
	checkScheme(uri string) error
}

type writerError struct {
	uri string
}

func (e *writerError) Error() string {
	return fmt.Sprintf("Error writing %s", e.uri)
}

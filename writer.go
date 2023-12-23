package main

import "fmt"

type Writer interface {
	write(bytes []byte, uri string) error
	checkScheme(uri string) error
}

type writerError struct {
	uri string
}

func (e *writerError) Error() string {
	return fmt.Sprintf("Error writing %s", e.uri)
}

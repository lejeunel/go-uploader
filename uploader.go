package main

import (
	"errors"
)

type ReadWriter struct {
	reader Reader
	writer Writer
}

type uploadError struct {
	msg            string
	uriSource      string
	uriDestination string
}

func (u ReadWriter) transfer(input string, output string) error {
	bytes, err_read := u.reader.read(input)
	err_write := u.writer.write(bytes, output)

	return errors.Join(err_read, err_write)
}

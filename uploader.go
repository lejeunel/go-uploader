package main

import (
	"errors"
)

type uploader struct {
	reader
	writer
}

type uploadError struct {
	msg    string
	uriIn  string
	uriOut string
}

func (u uploader) transfer(input string, output string) error {
	bytes, err_read := u.read(input)
	err_write := u.write(bytes, output)

	return errors.Join(err_read, err_write)
}

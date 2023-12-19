package main

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

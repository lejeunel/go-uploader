package main

import (
	"fmt"
	"strings"
)

type mockReader struct {
	writer writer
}

type mockWriter struct {
	reader reader
}

func (r *mockReader) read(uri string) []byte {
	return []byte{4, 2}
}

func (r *mockReader) check_scheme(uri string) error {
	if strings.HasPrefix(uri, "file:///") {
		return nil
	} else {
		return &schemeError{provided_uri: uri, allowed_scheme: uri}

	}
}

func (r *mockWriter) check_scheme(uri string) error {
	if strings.HasPrefix(uri, "scheme://") {
		return nil
	} else {
		return &schemeError{provided_uri: uri, allowed_scheme: uri}

	}
}

func (r *mockReader) scan(uri string) ([]string, error) {
	check_scheme_error := r.check_scheme(uri)
	if check_scheme_error != nil {
		return make([]string, 0), check_scheme_error
	}

	root := "file:///path/to/dir/"
	files := make([]string, 2)
	for i := 0; i < 2; i++ {
		files[i] = root + fmt.Sprintf("file_%03d.ext", i)
	}
	if uri == root {
		return files, nil
	} else {
		return make([]string, 0), nil
	}

}

func (w *mockWriter) write(bytes []byte, uri string) bool {
	return true
}

func NewMockUploader() *uploader {
	return &uploader{reader: &mockReader{},
		writer: &mockWriter{}}
}

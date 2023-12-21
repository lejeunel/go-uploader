package main

import (
	"fmt"
	"strings"
)

type mockReader struct {
	reader   reader
	dataPath string
	nFiles   int
}

type mockWriter struct {
	writer writer
}

func (r *mockReader) read(uri string) ([]byte, error) {
	return []byte{4, 2}, nil
}

func (r *mockReader) checkScheme(uri string) error {
	if strings.HasPrefix(uri, "file:///") {
		return nil
	} else {
		return &schemeError{provided_uri: uri, allowed_scheme: uri}

	}
}

func (r *mockReader) checkExists(uri string) error {
	if uri == r.dataPath {
		return nil
	} else {
		return &sourceError{provided_uri: uri}

	}
}

func (r *mockWriter) checkScheme(uri string) error {
	if strings.HasPrefix(uri, "scheme://") {
		return nil
	} else {
		return &schemeError{provided_uri: uri, allowed_scheme: uri}

	}
}

func (r *mockReader) scan(uri string) ([]string, error) {

	files := make([]string, r.nFiles)
	for i := 0; i < r.nFiles; i++ {
		files[i] = r.dataPath + fmt.Sprintf("file_%05d.ext", i)
	}
	return files, nil

}

func (w *mockWriter) write(bytes []byte, uri string) error {
	return nil
}

func NewMockUploader(nFiles int) *uploader {

	return &uploader{reader: &mockReader{dataPath: "file:///path/to/data/", nFiles: nFiles},
		writer: &mockWriter{}}
}

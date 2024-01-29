package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type MockReader struct {
	reader   Reader
	dataPath string
	nFiles   int
}

type MockWriter struct {
	writer Writer
}

func (r *MockReader) read(uri string) ([]byte, error) {
	return []byte{4, 2}, nil
}

func (r *MockReader) checkScheme(uri string) error {
	if strings.HasPrefix(uri, "file:///") {
		return nil
	} else {
		return &schemeError{provided_uri: uri, allowed_scheme: uri}

	}
}

func (r *MockReader) checkExists(uri string) error {
	if uri == r.dataPath {
		return nil
	} else {
		return &sourceError{provided_uri: uri}

	}
}

func (r *MockWriter) checkScheme(uri string) error {
	if strings.HasPrefix(uri, "scheme://") {
		return nil
	} else {
		return &schemeError{provided_uri: uri, allowed_scheme: uri}

	}
}

func (r *MockReader) scan(uri string) ([]string, error) {

	files := make([]string, r.nFiles)
	for i := 0; i < r.nFiles; i++ {
		files[i] = r.dataPath + fmt.Sprintf("file_%05d.ext", i)
	}
	return files, nil

}

func (w *MockWriter) write(bytes []byte, uri string) error {
	return nil
}

func NewMockReadWriter(nFiles int) *ReadWriter {

	return &ReadWriter{reader: &MockReader{dataPath: "file:///path/to/data/", nFiles: nFiles},
		writer: &MockWriter{}}
}

func NewMockJobManager() *jobManager {
	logger := MakeLogger(log.WarnLevel)
	return &jobManager{readWriter: *NewMockReadWriter(4), Store: NewMockStore(),
		logger:   logger,
		nWorkers: 2}
}

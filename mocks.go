package main

import (
	"fmt"
	"strings"
)

type mockReader struct {
	reader   reader
	dataPath string
}

type mockWriter struct {
	writer writer
}

type mockStore struct {
	store store
	jobs  []*job
}

func (r *mockReader) read(uri string) []byte {
	return []byte{4, 2}
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

func (r *mockReader) scan(uri string) []string {

	files := make([]string, 2)
	for i := 0; i < 2; i++ {
		files[i] = r.dataPath + fmt.Sprintf("file_%03d.ext", i)
	}
	return files

}

func (w *mockWriter) write(bytes []byte, uri string) bool {
	return true
}

func (s mockStore) commit(job *job) bool {
	return true
}

func (s mockStore) getJob(uriSource string, uriDestination string) (*job, error) {
	for _, job := range s.jobs {
		if (job.uriSource == uriSource) && (job.uriDestination == uriDestination) {
			return job, nil
		}
	}
	return nil, &jobNotFoundError{}
}

func NewMockUploader() *uploader {
	return &uploader{reader: &mockReader{dataPath: "file:///path/to/data/"},
		writer: &mockWriter{}}
}

func NewMockStore() *mockStore {
	return &mockStore{}
}

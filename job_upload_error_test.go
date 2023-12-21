package main

import (
	"errors"
	"testing"
)

type mockReaderErrorOnRead struct {
	mockReader
}

func (r *mockReaderErrorOnRead) read(uri string) ([]byte, error) {
	return []byte{}, &readerError{uri: uri}
}

type mockWriterErrorOnWrite struct {
	mockWriter
}

func (r *mockWriterErrorOnWrite) write(bytes []byte, uri string) error {
	return &writerError{uri: uri}
}

func TestJobUploadReadError(t *testing.T) {
	uploader := uploader{reader: &mockReaderErrorOnRead{mockReader{dataPath: "file:///path/to/data/",
		nFiles: 10}},
		writer: &mockWriter{}}
	jm := &jobManager{uploader: uploader, store: NewMockStore(), nWorkers: 5}
	_, err := MakeCompletedJob(jm)

	var got *readerError
	isReadError := errors.As(err, &got)

	if !isReadError {
		t.Fatalf("expected to retrieve reader error but got %v", err)
	}
}

func TestJobUploadWriteError(t *testing.T) {
	uploader := uploader{reader: &mockReader{dataPath: "file:///path/to/data/",
		nFiles: 10},
		writer: &mockWriterErrorOnWrite{}}
	jm := &jobManager{uploader: uploader, store: NewMockStore(), nWorkers: 5}
	_, err := MakeCompletedJob(jm)

	var got *writerError
	isWriteError := errors.As(err, &got)

	if !isWriteError {
		t.Fatalf("expected to retrieve writer error but got %v", err)
	}
}

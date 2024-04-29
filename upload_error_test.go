package main

import (
	"errors"
	"testing"
)

type MockReaderErrorOnRead struct {
	MockReader
}

func (r *MockReaderErrorOnRead) read(uri string) ([]byte, error) {
	return []byte{}, &readerError{uri: uri}
}

type mockWriterErrorOnWrite struct {
	MockWriter
}

func (r *mockWriterErrorOnWrite) write(bytes []byte, uri string) error {
	return &writerError{uri: uri}
}

func TestUploadReadError(t *testing.T) {
	readWriter := ReadWriter{reader: &MockReaderErrorOnRead{MockReader{dataPath: "file:///path/to/data/",
		nFiles: 10}},
		writer: &MockWriter{}}
	jm := &jobManager{readWriter: readWriter, store: NewMockStore(), nWorkers: 5}
	_, err := MakeCompletedJob(jm)

	var got *readerError
	isReadError := errors.As(err, &got)

	if !isReadError {
		t.Fatalf("expected to retrieve reader error but got %v", err)
	}
}

func TestUploadWriteError(t *testing.T) {
	readWriter := ReadWriter{reader: &MockReader{dataPath: "file:///path/to/data/",
		nFiles: 10},
		writer: &mockWriterErrorOnWrite{}}
	jm := &jobManager{readWriter: readWriter, store: NewMockStore(), nWorkers: 5}
	_, err := MakeCompletedJob(jm)

	var got *writerError
	isWriteError := errors.As(err, &got)

	if !isWriteError {
		t.Fatalf("expected to retrieve writer error but got %v", err)
	}
}

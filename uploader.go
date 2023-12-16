package main

type Reader interface {
	read() bool
	scan() []string
}

type Writer interface {
	write() bool
}

type FSReader struct {
	fs fileSystem
}

func (FSReader) scan(uri string) {

}

type Uploader interface {
	upload(input string, output string) bool
}

type MockUploader struct {
	uploader Uploader
	reader   Reader
	writer   Writer
}

func (MockUploader) upload(input string, output string) bool {
	return false
}

func NewMockUploader() *MockUploader {
	return &MockUploader{}
}

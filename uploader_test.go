package main

import (
	"testing"
)

func TestUploader(t *testing.T) {
	uploader := NewMockUploader()
	success := uploader.upload("input", "output")
	if !success {
		t.Fatalf(`Upload failed`)
	}
}

package main

import (
	"errors"
	"testing"
)

func TestRetrieveJobFromStore(t *testing.T) {
	jm := NewJobManager(*NewMockUploader(), *NewMockStore())
	job := MakeCompletedJob(jm)

	_, err := jm.getJob(job.uriSource, job.uriDestination)

	var got *jobNotFoundError
	isJobNotFoundError := errors.As(err, &got)

	if !isJobNotFoundError {
		t.Fatalf("expected to retrieve job but got none")
	}

}

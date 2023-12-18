package main

import (
	"errors"
	"testing"
)

func TestRetrieveJobFromStore(t *testing.T) {
	jm := NewJobManager(*NewMockUploader(), NewMockStore())
	job := MakeCompletedJob(jm)

	retrieved_job, err := jm.GetJob(job.uriSource, job.uriDestination)

	var got *jobNotFoundError
	isJobNotFoundError := errors.As(err, &got)

	if isJobNotFoundError {
		t.Fatalf("expected to retrieve job but got none")
	}

	if retrieved_job.status != done {
		t.Fatalf("expected to retrieve finished job but got status %v", retrieved_job.status)
	}

}

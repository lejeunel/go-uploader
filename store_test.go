package main

import (
	"errors"
	"testing"
)

func TestRetrieveCompletedFromStore(t *testing.T) {

	jm := NewMockJobManager()
	job, _ := MakeCompletedJob(jm)

	retrieved_job, err := jm.store.FindJob(job.UriSource, job.UriDestination)

	var got *jobNotFoundError
	isJobNotFoundError := errors.As(err, &got)

	if isJobNotFoundError || (err != nil) {
		t.Fatalf("expected to retrieve job but got none. Error %v", err)
	}

	if (retrieved_job.Status != done) || (err != nil) {
		t.Fatalf("expected to retrieve finished job but got status %v. Error %v", retrieved_job.Status, err)
	}

}

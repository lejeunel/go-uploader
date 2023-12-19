package main

import (
	"errors"
	"testing"
)

func TestRetrieveCompletedJobFromStore(t *testing.T) {

	jm := NewMockJobManager()
	job := MakeCompletedJob(jm)

	retrieved_job, err := jm.GetJob(job.UriSource, job.UriDestination)

	var got *jobNotFoundError
	isJobNotFoundError := errors.As(err, &got)

	if isJobNotFoundError || (err != nil) {
		t.Fatalf("expected to retrieve job but got none. Error %v", err)
	}

	if (retrieved_job.Status != done) || (err != nil) {
		t.Fatalf("expected to retrieve finished job but got status %v. Error %v", retrieved_job.Status, err)
	}

}

func TestRetrieveCreatedJobFromStore(t *testing.T) {

	jm := NewMockJobManager()
	job, _ := jm.CreateJob("file:///path/to/data/", "scheme://path/to/data/")
	job, err := jm.GetJob(job.UriSource, job.UriDestination)

	jm.ParseJob(job)
	jm.TransferJob(job)

	var got *jobNotFoundError
	isJobNotFoundError := errors.As(err, &got)

	if isJobNotFoundError || (err != nil) {
		t.Fatalf("expected to retrieve job but got none. Error %v", err)
	}

	if (job.Status != done) || (err != nil) {
		t.Fatalf("expected to retrieve finished job but got status %v. Error %v", job.Status, err)
	}

}

func TestRetrieveParsedJobFromStore(t *testing.T) {

	jm := NewMockJobManager()
	job, _ := jm.CreateJob("file:///path/to/data/", "scheme://path/to/data/")
	jm.ParseJob(job)

	job, err := jm.GetJob(job.UriSource, job.UriDestination)

	jm.TransferJob(job)

	var got *jobNotFoundError
	isJobNotFoundError := errors.As(err, &got)

	if isJobNotFoundError || (err != nil) {
		t.Fatalf("expected to retrieve job but got none. Error %v", err)
	}

	if (job.Status != done) || (err != nil) {
		t.Fatalf("expected to retrieve finished job but got status %v. Error %v", job.Status, err)
	}

}

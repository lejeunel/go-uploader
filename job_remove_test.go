package main

import (
	"errors"
	"testing"
)

func TestDelete(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := MakeCompletedJob(jm)
	id := job.ID

	err := jm.store.DeleteJob(job)
	if err != nil {
		t.Fatalf("expected to delete job without error. Got %v", err)
	}

	_, err_find := jm.FindJob(job.UriSource, job.UriDestination)

	var got *jobNotFoundError
	isJobNotFoundError := errors.As(err_find, &got)

	if !isJobNotFoundError {
		t.Fatalf("expected to get a job-not-found. Got %v", err_find)
	}

	transactions, _ := jm.store.getTransactions(id)
	if len(transactions) != 0 {
		t.Fatalf("expected to retrieve no transactions. Got %v", len(transactions))

	}

}

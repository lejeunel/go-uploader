package main

import (
	"testing"
)

func TestJobParse(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := jm.CreateJob("file:///path/to/data/", "scheme://path/to/data/")
	jm.ParseJob(job)
	if job.Status != parsed {
		t.Fatalf("expected job status %v, got %T", parsed, job.Status)
	}

	transactions := job.Transactions
	if len(transactions) < 1 {
		t.Fatalf("expected transactions, got none")
	}
}

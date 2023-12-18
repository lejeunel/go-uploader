package main

import (
	"testing"
)

func TestJobParse(t *testing.T) {
	jm := NewJobManager(*NewMockUploader(), NewMockStore())
	job, _ := jm.create("file:///path/to/data/", "scheme://path/to/data/")
	jm.parse(job)
	if job.status != parsed {
		t.Fatalf("expected job status %v, got %T", parsed, job.status)
	}

	transactions := job.transactions
	if len(transactions) < 1 {
		t.Fatalf("expected transactions, got none")
	}
}

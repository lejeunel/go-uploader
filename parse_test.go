package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := jm.Create("file:///path/to/data/", "scheme://path/to/data/")
	jm.Parse(job)
	if job.Status != parsed {
		t.Fatalf("expected job status %v, got %T", parsed, job.Status)
	}

	transactions := job.Transactions
	if len(transactions) < 1 {
		t.Fatalf("expected transactions, got none")
	}
}

func TestCannotReParse(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := MakeCompletedJob(jm)

	timestamp := job.Transactions[0].CreatedAt

	rejob, _ := jm.Parse(job)
	new_timestamp := rejob.Transactions[0].CreatedAt

	if timestamp != new_timestamp {
		t.Fatalf("expected re-parsing to skip adding transactions. Got timestamps %s and %s", timestamp,
			new_timestamp)
	}

}

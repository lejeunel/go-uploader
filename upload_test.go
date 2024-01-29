package main

import (
	"testing"
)

func MakeCompletedJob(jm *jobManager) (*Job, error) {
	job, err := jm.Run("file:///path/to/data/", "scheme://path/to/data/")
	if err != nil {
		return nil, err
	}

	return job, nil
}

func TestUpload(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := MakeCompletedJob(jm)
	if job.Status != done {
		t.Fatalf("expected job status %v, got %T", done, job.Status)
	}
	done_transactions := job.numDoneTransactions()
	n_transactions := job.numTransactions()
	if done_transactions != len(job.Transactions) {
		t.Fatalf("expected all of %v transactions to complete, got %v", n_transactions, done_transactions)
	}
}

func TestResumeAtTransfer(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := MakeCompletedJob(jm)
	job.Transactions[0].Status = pending
	job.Status = parsed

	job_, _ := jm.Transfer(job)
	done_transactions := job_.numDoneTransactions()
	if done_transactions != len(job_.Transactions) {
		t.Fatalf("expected all transactions to complete, got %v", done_transactions)
	}
}

func TestResumeAtParsing(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := jm.Create("file:///path/to/data/", "scheme://path/to/data/")
	jm.Parse(job)

	new_job, _ := jm.Resume(job)
	done_transactions := new_job.numDoneTransactions()
	if done_transactions != len(new_job.Transactions) {
		t.Fatalf("expected all transactions to complete, got %v", done_transactions)
	}
}

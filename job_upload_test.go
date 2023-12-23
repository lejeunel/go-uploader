package main

import (
	"testing"
)

func MakeCompletedJob(jm *jobManager) (*Job, error) {
	job, err_create := jm.CreateJob("file:///path/to/data/", "scheme://path/to/data/")
	if err_create != nil {
		return &Job{}, err_create
	}
	job, err_parse := jm.ParseJob(job)
	if err_parse != nil {
		return &Job{}, err_parse
	}
	job, err_trf := jm.TransferJob(job)
	if err_trf != nil {
		return &Job{}, err_trf
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

func TestResumedUpload(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := MakeCompletedJob(jm)
	job.Transactions[0].Status = pending
	job.Status = parsed

	job_, _ := jm.TransferJob(job)
	done_transactions := job_.numDoneTransactions()
	if done_transactions != len(job_.Transactions) {
		t.Fatalf("expected all transactions to complete, got %v", done_transactions)
	}
}

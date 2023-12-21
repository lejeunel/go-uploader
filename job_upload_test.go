package main

import (
	"errors"
	"testing"
)

func MakeCompletedJob(jm *jobManager) (*Job, error) {
	job, err_create := jm.CreateJob("file:///path/to/data/", "scheme://path/to/data/")
	err_parse := jm.ParseJob(job)
	_, err_trf := jm.TransferJob(job)

	return job, errors.Join(err_create, err_parse, err_trf)
}

func TestJobUpload(t *testing.T) {
	jm := NewMockJobManager()
	job, _ := MakeCompletedJob(jm)
	if job.Status != done {
		t.Fatalf("expected job status %v, got %T", done, job.Status)
	}
}

func TestJobResumedUpload(t *testing.T) {
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

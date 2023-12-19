package main

import (
	"testing"
)

func MakeCompletedJob(jm *jobManager) *Job {
	job, _ := jm.CreateJob("file:///path/to/data/", "scheme://path/to/data/")
	jm.ParseJob(job)
	jm.TransferJob(job)

	return job
}

func TestJobUpload(t *testing.T) {
	jm := NewMockJobManager()
	job := MakeCompletedJob(jm)
	if job.Status != done {
		t.Fatalf("expected job status %v, got %T", done, job.Status)
	}
}

func TestJobResumedUpload(t *testing.T) {
	jm := NewMockJobManager()
	job := MakeCompletedJob(jm)
	job.Transactions[0].Status = pending
	job.Status = parsed

	job_ := jm.TransferJob(job)
	done_transactions := job_.numDoneTransactions()
	if done_transactions != len(job_.Transactions) {
		t.Fatalf("expected all transactions to complete, got %v", done_transactions)
	}
}

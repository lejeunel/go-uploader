package main

import (
	"testing"
)

func MakeCompletedJob(jm *jobManager) *Job {
	job, _ := jm.create("file:///path/to/data/", "scheme://path/to/data/")
	jm.parse(job)
	jm.upload(job)

	return job
}

func TestJobUpload(t *testing.T) {
	jm := NewJobManager(*NewMockUploader(), NewMockStore(), 2)
	job := MakeCompletedJob(jm)
	if job.Status != done {
		t.Fatalf("expected job status %v, got %T", done, job.Status)
	}
}

func TestJobResumedUpload(t *testing.T) {
	jm := NewJobManager(*NewMockUploader(), NewMockStore(), 2)
	job := MakeCompletedJob(jm)
	job.Transactions[0].status = pending
	job.Status = parsed
	numTransactionsUploaded := jm.upload(job)
	if numTransactionsUploaded > 1 {
		t.Fatalf("expected one transaction to complete, got %T", numTransactionsUploaded)
	}
}

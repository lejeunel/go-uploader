package main

import (
	"testing"
)

func MakeCompletedJob(jm *jobManager) *job {
	job, _ := jm.create("file:///path/to/data/", "scheme://path/to/data/")
	jm.parse(job)
	jm.upload(job)

	return job
}

func TestJobUpload(t *testing.T) {
	jm := NewJobManager(*NewMockUploader(), *NewMockStore())
	job := MakeCompletedJob(jm)
	if job.status != done {
		t.Fatalf("expected job status %v, got %T", done, job.status)
	}
}

func TestJobResumedUpload(t *testing.T) {
	jm := NewJobManager(*NewMockUploader(), *NewMockStore())
	job := MakeCompletedJob(jm)
	job.transactions[0].status = pending
	job.status = parsed
	numTransactionsUploaded := jm.upload(job)
	if numTransactionsUploaded > 1 {
		t.Fatalf("expected one transaction to complete, got %T", numTransactionsUploaded)
	}
}

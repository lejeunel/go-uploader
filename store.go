package main

type store interface {
	commit(job *job) bool
	getJob(uriSource string, uriDestination string) (*job, error)
}

type jobNotFoundError struct {
}

func (e *jobNotFoundError) Error() string {
	return "Job not found error"
}

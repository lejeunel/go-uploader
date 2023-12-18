package main

type store interface {
	Commit(job *job)
	GetJob(uriSource string, uriDestination string) (*job, error)
}

type jobNotFoundError struct {
}

func (e *jobNotFoundError) Error() string {
	return "Job not found error"
}

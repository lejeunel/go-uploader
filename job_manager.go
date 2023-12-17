package main

type jobInitError struct {
}

func (e *jobInitError) Error() string {
	return "Error initializing job"
}

type jobManager struct {
	uploader
}

func (m *jobManager) create(uri_source string, uri_destination string) (*job, error) {
	_, e := m.uploader.scan(uri_source)

	if e != nil {
		return &job{}, e
	}

	e = m.uploader.writer.check_scheme(uri_destination)
	if e != nil {
		return &job{}, e
	}

	return &job{}, nil
}

func NewJobManager(uploader uploader) *jobManager {
	return &jobManager{uploader: uploader}
}

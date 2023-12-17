package main

import (
	"time"
)

type jobInitError struct {
}

func (e *jobInitError) Error() string {
	return "Error initializing job"
}

type jobManager struct {
	uploader
}

func (m *jobManager) create(uriSource string, uriDestination string) (*job, error) {
	eSource := m.uploader.reader.checkScheme(uriSource)

	if eSource != nil {
		return nil, eSource
	}

	eDestination := m.uploader.writer.checkScheme(uriDestination)
	if eDestination != nil {
		return nil, eDestination
	}

	eSourceExists := m.uploader.reader.checkExists(uriSource)
	if eSourceExists != nil {
		return nil, eSourceExists
	}

	return &job{uriSource: uriSource, uriDestination: uriDestination,
		lastStatus: created, createdAt: time.Now(), updatedAt: time.Now()}, nil
}

func NewJobManager(uploader uploader) *jobManager {
	return &jobManager{uploader: uploader}
}

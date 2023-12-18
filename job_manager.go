package main

import (
	"strings"
	"time"
)

type jobInitError struct {
}

func (e *jobInitError) Error() string {
	return "Error initializing job"
}

type jobManager struct {
	uploader
	store
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

	job := &job{uriSource: uriSource, uriDestination: uriDestination,
		status: created, createdAt: time.Now(), updatedAt: time.Now()}
	m.Commit(job)

	return job, nil
}

func (m *jobManager) parse(job *job) {
	inURIs := m.reader.scan(job.uriSource)
	var transactions []*transaction

	for _, f := range inURIs {
		parts := strings.Split(f, "/")
		stem := parts[len(parts)-1]
		t := transaction{uriIn: f, uriOut: job.uriDestination + stem}
		transactions = append(transactions, &t)

	}
	job.transactions = transactions
	job.status = parsed

}

func (m *jobManager) upload(job *job) int {
	numTransactionsUploaded := 0
	for _, t := range job.transactions {
		if t.status == pending {
			bytes := m.uploader.read(t.uriIn)
			m.uploader.write(bytes, t.uriOut)
			t.status = transferred
			numTransactionsUploaded++
		}
	}

	job.status = done

	return numTransactionsUploaded
}

func NewJobManager(uploader uploader, store store) *jobManager {
	return &jobManager{uploader: uploader, store: store}
}

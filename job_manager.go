package main

import (
	"github.com/google/uuid"
	"strings"
	"sync"
)

type jobManager struct {
	uploader
	store
	nWorkers int
}

func (m *jobManager) CreateJob(uriSource string, uriDestination string) (*Job, error) {
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

	job := &Job{UriSource: uriSource, UriDestination: uriDestination,
		Status: created}
	m.store.AppendJob(job)

	return job, nil
}

func (m *jobManager) ParseJob(job *Job) {
	inURIs := m.reader.scan(job.UriSource)
	var transactions []Transaction

	for _, f := range inURIs {
		parts := strings.Split(f, "/")
		stem := parts[len(parts)-1]
		t := Transaction{ID: uuid.New(), uriIn: f, uriOut: job.UriDestination + stem}
		transactions = append(transactions, t)
	}
	job.Transactions = transactions
	job.Status = parsed
	m.store.UpdateJob(job)

}

func (m *jobManager) transferWorker(worker_id int, transactions <-chan Transaction,
	results chan<- Transaction,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for transaction := range transactions {
		bytes := m.uploader.read(transaction.uriIn)
		m.uploader.write(bytes, transaction.uriOut)
		transaction.Status = transferred
		results <- transaction
	}
}

func (m *jobManager) updateTransactionWorker(results <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for transaction := range results {
		m.store.UpdateTransaction(&transaction)
	}
}

func (m *jobManager) TransferJob(job *Job) *Job {
	var pending_transactions []Transaction
	for _, t := range job.Transactions {
		if t.Status == pending {
			pending_transactions = append(pending_transactions, t)
		}
	}

	todo := make(chan Transaction, len(pending_transactions))
	results := make(chan Transaction, len(pending_transactions))
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	for w := 0; w < m.nWorkers; w++ {
		wg.Add(1)
		go m.transferWorker(w, todo, results, &wg)
	}

	go m.updateTransactionWorker(results, &wg2)

	for _, t := range pending_transactions {
		todo <- t
	}
	close(todo)

	wg.Wait()
	wg2.Wait()

	job.Status = done
	done_job := m.store.UpdateJob(job)

	return done_job

}

func NewJobManager(uploader uploader, store store, nWorkers int) *jobManager {
	return &jobManager{uploader: uploader, store: store, nWorkers: nWorkers}
}

func NewMockJobManager() *jobManager {
	return &jobManager{uploader: *NewMockUploader(1000), store: NewMockStore(), nWorkers: 10}
}

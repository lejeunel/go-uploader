package main

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
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

func (m *jobManager) ParseJob(job *Job) error {
	inURIs, err := m.reader.scan(job.UriSource)
	if err != nil {
		return err
	}

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

	return nil

}

func (m *jobManager) updateTransactionWorker(results <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for transaction := range results {
		m.store.UpdateTransaction(&transaction)
	}
}

func (m *jobManager) transferWorker(ctx context.Context, worker_id int, transactions <-chan Transaction,
	results chan<- Transaction) error {
	for transaction := range transactions {
		bytes, err_read := m.uploader.read(transaction.uriIn)
		err_write := m.uploader.write(bytes, transaction.uriOut)
		transaction.Status = transferred
		results <- transaction
		joined_err := errors.Join(err_read, err_write)
		if joined_err != nil {
			return joined_err
		}
	}
	return nil
}

func (m *jobManager) TransferJob(job *Job) (*Job, error) {
	var pending_transactions []Transaction
	for _, t := range job.Transactions {
		if t.Status == pending {
			pending_transactions = append(pending_transactions, t)
		}
	}

	todo := make(chan Transaction, len(pending_transactions))
	results := make(chan Transaction, len(pending_transactions))
	ctx := context.Background()
	eg, egCtx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	for w := 0; w < m.nWorkers; w++ {
		eg.Go(func() error {
			return m.transferWorker(egCtx, w, todo, results)
		})
	}

	for _, t := range pending_transactions {
		todo <- t
	}
	close(todo)

	go m.updateTransactionWorker(results, &wg)

	if err := eg.Wait(); err != nil {
		return job, err
	}

	wg.Wait()

	job.Status = done
	done_job := m.store.UpdateJob(job)

	return done_job, nil

}

func NewJobManager(uploader uploader, store store, nWorkers int) *jobManager {
	return &jobManager{uploader: uploader, store: store, nWorkers: nWorkers}
}

func NewMockJobManager() *jobManager {
	return &jobManager{uploader: *NewMockUploader(10), store: NewMockStore(), nWorkers: 5}
}

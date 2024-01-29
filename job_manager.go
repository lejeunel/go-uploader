package main

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
)

type jobManager struct {
	Store
	readWriter ReadWriter
	logger     *log.Logger
	nWorkers   int
}

func (m *jobManager) Run(uriSource string, uriDestination string) (*Job, error) {
	job, err_create := m.Create("file:///path/to/data/", "scheme://path/to/data/")
	if err_create != nil {
		return nil, err_create
	}
	job, err_parse := m.Parse(job)
	if err_parse != nil {
		return job, err_parse
	}
	job, err_trf := m.Transfer(job)
	if err_trf != nil {
		return job, err_trf
	}

	return job, nil
}

func (m *jobManager) Resume(job *Job) (*Job, error) {
	job, err_parse := m.Parse(job)
	if err_parse != nil {
		return job, err_parse
	}
	job, err_trf := m.Transfer(job)
	if err_trf != nil {
		return job, err_trf
	}

	return job, nil
}

func (m *jobManager) Create(uriSource string, uriDestination string) (*Job, error) {
	eSource := m.readWriter.reader.checkScheme(uriSource)
	eDestination := m.readWriter.writer.checkScheme(uriDestination)
	eSourceExists := m.readWriter.reader.checkExists(uriSource)
	duplicate_job, eJobNotFound := m.FindJob(uriSource, uriDestination)

	joined_err := errors.Join(eDestination, eSourceExists, eSource)
	if eJobNotFound == nil {
		return duplicate_job, errors.Join(joined_err, &duplicateJobError{job: duplicate_job})
	}

	if joined_err != nil {
		return nil, joined_err
	}

	job := &Job{UriSource: uriSource, UriDestination: uriDestination,
		Status: created}
	m.AppendJob(job)

	return job, nil
}

func (m *jobManager) Parse(job *Job) (*Job, error) {
	if job.Status >= parsed {
		return job, nil
	}

	inURIs, err := m.readWriter.reader.scan(job.UriSource)
	if err != nil {
		return job, err
	}

	var transactions []Transaction

	for _, f := range inURIs {
		parts := strings.Split(f, "/")
		stem := parts[len(parts)-1]
		t := Transaction{UriSource: f, UriDestination: job.UriDestination + stem}
		transactions = append(transactions, t)
	}
	job.Transactions = transactions
	job.Status = parsed
	err_upd := m.UpdateJobStatus(job)
	job, err_apd := m.AppendJobTransactions(job)

	return job, errors.Join(err_upd, err_apd)

}

func (m *jobManager) updateTransactionWorker(results <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for transaction := range results {
		m.UpdateTransactionStatus(&transaction)
		m.logger.WithFields(log.Fields{
			"in":  transaction.UriSource,
			"out": transaction.UriDestination}).Info("transferred")
	}
}

func (m *jobManager) transferWorker(ctx context.Context, worker_id int, transactions <-chan Transaction,
	results chan<- Transaction) error {
	for transaction := range transactions {
		bytes, err_read := m.readWriter.reader.read(transaction.UriSource)
		err_write := m.readWriter.writer.write(bytes, transaction.UriDestination)
		joined_err := errors.Join(err_read, err_write)
		if joined_err != nil {
			return joined_err
		} else {
			transaction.Status = transferred
			results <- transaction

		}
	}
	return nil
}

func (m *jobManager) Transfer(job *Job) (*Job, error) {
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

	for w := 0; w < m.nWorkers; w++ {
		eg.Go(func() error {
			return m.transferWorker(egCtx, w, todo, results)
		})
	}

	for _, t := range pending_transactions {
		todo <- t
	}
	close(todo)

	var wg sync.WaitGroup
	wg.Add(1)
	go m.updateTransactionWorker(results, &wg)

	err := eg.Wait()
	if err != nil {
		return job, err
	}
	close(results)

	wg.Wait()

	job.Status = done
	err_upd := m.UpdateJobStatus(job)
	done_job, err_find := m.FindJob(job.UriSource, job.UriDestination)

	return done_job, errors.Join(err, err_upd, err_find)

}

func NewJobManager(readWriter ReadWriter, store Store, nWorkers int) *jobManager {
	return &jobManager{readWriter: readWriter, Store: store, nWorkers: nWorkers}
}

func MakeLogger(level log.Level) *log.Logger {
	logger := log.New()
	logger.SetLevel(log.WarnLevel)
	logger.SetFormatter(&log.JSONFormatter{})
	return logger
}

package main

import (
	// "fmt"
	"strings"
	"sync"
)

type jobManager struct {
	uploader
	store
	nWorkers int
}

func (m *jobManager) create(uriSource string, uriDestination string) (*Job, error) {
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

func (m *jobManager) parse(job *Job) {
	inURIs := m.reader.scan(job.UriSource)
	var transactions []Transaction

	for _, f := range inURIs {
		parts := strings.Split(f, "/")
		stem := parts[len(parts)-1]
		// fmt.Println("parsed", f, "->", job.UriDestination+stem)
		t := Transaction{uriIn: f, uriOut: job.UriDestination + stem}
		transactions = append(transactions, t)
	}
	job.Transactions = transactions
	job.Status = parsed
	m.store.UpdateJob(job)

}

func (m *jobManager) upload_worker(worker_id int, transactions <-chan Transaction,
	results chan<- Transaction,
	wg *sync.WaitGroup) {
	defer wg.Done()
	for transaction := range transactions {
		// fmt.Println("Worker", worker_id, transaction.uriIn, "->", transaction.uriOut)
		bytes := m.uploader.read(transaction.uriIn)
		m.uploader.write(bytes, transaction.uriOut)
		transaction.status = transferred
		results <- transaction
	}
}

func (m *jobManager) results_worker(results <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for transaction := range results {
		// fmt.Println("Writing result", transaction.uriIn, "->", transaction.uriOut, "/ status", transaction.status)
		m.store.UpdateTransaction(&transaction)
	}
}

func (m *jobManager) upload(job *Job) int {
	var pending_transactions []Transaction
	for _, t := range job.Transactions {
		if t.status == pending {
			pending_transactions = append(pending_transactions, t)
		}
	}

	nWorkers := 4
	todo := make(chan Transaction, len(pending_transactions))
	results := make(chan Transaction, len(pending_transactions))
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup

	// fmt.Println("Spawning", nWorkers, "workers")
	for w := 0; w < nWorkers; w++ {
		wg.Add(1)
		go m.upload_worker(w, todo, results, &wg)
	}

	// fmt.Println("Spawning results worker")
	go m.results_worker(results, &wg2)

	// fmt.Println("Appending to transaction channel...")
	for _, t := range pending_transactions {
		todo <- t
	}
	close(todo)

	wg.Wait()
	wg2.Wait()

	job.Status = done
	m.store.UpdateJob(job)

	return 0

}

func NewJobManager(uploader uploader, store store, nWorkers int) *jobManager {
	return &jobManager{uploader: uploader, store: store, nWorkers: nWorkers}
}

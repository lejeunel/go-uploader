package main

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type store interface {
	AppendJob(job *Job)
	UpdateJob(job *Job)
	UpdateTransaction(t *Transaction)
	GetJob(uriSource string, uriDestination string) (*Job, error)
}

type SQLiteStore struct {
	store store
	db    *gorm.DB
}

func (s *SQLiteStore) AppendJob(job *Job) {
	job.ID = uuid.New()
	// fmt.Println("creating job with status/ID", job.Status, job.ID)
	s.db.Create(job)
}

func (s *SQLiteStore) UpdateJob(job *Job) {
	fmt.Println("updating job with status/ID", job.Status, job.ID)
	s.db.Save(&job)
}

func (s *SQLiteStore) UpdateTransaction(t *Transaction) {
	s.db.Save(t)
}

func (s *SQLiteStore) GetJob(uriSource string, uriDestination string) (*Job, error) {
	var job = Job{}
	result := s.db.Where("uri_source = ? AND uri_destination = ?", uriSource, uriDestination).First(&job)
	// fmt.Println("retrieved job with status/ID", job.Status, job.ID)
	return &job, result.Error
}

func NewStore(url string) *SQLiteStore {

	db, err := gorm.Open(sqlite.Open(url), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Job{}, &Transaction{})

	return &SQLiteStore{db: db}

}

func NewMockStore() *SQLiteStore {
	return NewStore("file:mock.sqlite?cache=shared&_journal_mode=WAL")
}

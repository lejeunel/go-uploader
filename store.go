package main

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type store interface {
	AppendJob(job *Job) *Job
	UpdateJob(job *Job) *Job
	UpdateTransaction(t *Transaction) *Transaction
	GetJob(uriSource string, uriDestination string) (*Job, error)
}

type SQLiteStore struct {
	store store
	db    *gorm.DB
}

func (s *SQLiteStore) AppendJob(job *Job) *Job {
	job.ID = uuid.New()
	s.db.Create(job)
	return job
}

func (s *SQLiteStore) UpdateJob(job *Job) *Job {
	s.db.Save(&job)
	updated_job, _ := s.GetJob(job.UriSource, job.UriDestination)
	return updated_job
}

func (s *SQLiteStore) UpdateTransaction(t *Transaction) *Transaction {
	fmt.Println("updating transaction with status/ID", t.Status, t.ID)
	s.db.Save(t)
	return t
}

func (s *SQLiteStore) GetJob(uriSource string, uriDestination string) (*Job, error) {
	var job = Job{}
	result := s.db.Where("uri_source = ? AND uri_destination = ?", uriSource, uriDestination).First(&job).Preload("Transactions")
	return &job, result.Error
}

func NewStore(db *gorm.DB) *SQLiteStore {

	// Migrate the schema
	db.AutoMigrate(&Job{}, &Transaction{})

	return &SQLiteStore{db: db}

}

func NewMockStore() *SQLiteStore {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(1)
	return NewStore(db)
}

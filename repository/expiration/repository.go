package expiration

import (
	"database/sql"

	"github.com/deadcheat/cashew"
)

// Repository is a struct will implement cashew.ExpirationRepository
type Repository struct {
	db *sql.DB
}

// New create new ExpirationRepository
func New(db *sql.DB) cashew.ExpirationRepository {
	return &Repository{db: db}
}

func (r *Repository) FindAll() ([]*cashew.Ticket, error) {

	return nil, nil
}

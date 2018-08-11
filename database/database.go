package database

import (
	"database/sql"
)

// Connector database connector interface
type Connector interface {
	Open() (*sql.DB, error)
}

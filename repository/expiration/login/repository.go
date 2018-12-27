package login

import (
	"database/sql"
	"time"

	"github.com/deadcheat/cashew/foundation"
	"github.com/deadcheat/cashew/timer"

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

// FindAll select all expired tickets
func (r *Repository) FindAll() (ts []*cashew.Ticket, err error) {

	var stmt *sql.Stmt

	stmt, err = r.db.Prepare(selectByTimeoutTicketQuery)
	if err != nil {
		return
	}
	defer stmt.Close()
	now := timer.Local.Now()
	hardTimeOut := now.Add(-1 * time.Second * time.Duration(foundation.App().GrantingHardTimeout))

	var rows *sql.Rows
	rows, err = stmt.Query(cashew.TicketTypeLogin.String(), hardTimeOut)
	if err != nil {
		return
	}
	defer rows.Close()
	ts = make([]*cashew.Ticket, 0)
	for rows.Next() {
		var (
			ticket  cashew.Ticket
			typeStr sql.NullString
		)
		err = rows.Scan(
			&ticket.ID,
			&typeStr,
			&ticket.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if typeStr.Valid {
			// NullString always return nil as error
			tmp, _ := typeStr.Value()
			ticket.Type = cashew.ParseTicketType(tmp.(string))
		}

		ts = append(ts, &ticket)
	}
	return ts, nil
}

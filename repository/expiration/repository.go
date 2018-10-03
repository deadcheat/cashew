package expiration

import (
	"database/sql"
	"fmt"
	"strings"
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

	inQuery := strings.Join([]string{cashew.TicketTypeTicketGranting.String(), cashew.TicketTypeProxyGranting.String()}, " ,")

	stmt, err = r.db.Prepare(fmt.Sprintf(selectByTimeoutTicketQuery, inQuery))
	if err != nil {
		return
	}
	defer stmt.Close()
	var rows *sql.Rows
	now := timer.Local.Now()
	hardTimeOut := now.Add(time.Second * time.Duration(foundation.App().GrantingHardTimeout))
	timeOut := now.Add(time.Second * time.Duration(foundation.App().GrantingDefaultExpire))
	rows, err = stmt.Query(timeOut, hardTimeOut)
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
			&ticket.LastReferencedAt,
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
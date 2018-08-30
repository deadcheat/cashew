package ticket

import (
	"database/sql"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/timer"
	"github.com/deadcheat/cashew/values/errs"
)

var (
	// updater for ticket_expires
	updateTicketLastReferenced ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		switch t.Type {
		case cashew.TicketTypeTicketGranting, cashew.TicketTypeProxyGranting, cashew.TicketTypeService, cashew.TicketTypeProxy:
			// do nothing
		default:
			return nil
		}
		// set automatically
		if t.LastReferencedAt == nil {
			now := timer.Local.Now()
			t.LastReferencedAt = &now
		}
		var stmt *sql.Stmt
		stmt, err := tx.Prepare(updateConsumeQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()
		var res sql.Result
		if res, err = stmt.Exec(timer.Local.Now(), t.ID); err != nil {
			return err
		}
		var count int64
		count, err = res.RowsAffected()
		if count == 0 {
			return errs.ErrNoTicketID
		}
		return err
	}
)

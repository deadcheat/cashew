package ticket

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/deadcheat/cashew"
)

// deleters
var deleteAccessors = []ticketAccessor{
	deleteAllRelatedTables,
	deleteTicket,
}

const (
	ticketExtraAttributes = "ticket_extra_attribute"
	ticketGrantTicket     = "ticket_grant_ticket"
	ticketIOU             = "ticket_iou"
	ticketLastReferenced  = "ticket_last_referenced"
	ticketService         = "ticket_service"
	ticketType            = "ticket_type"
	ticketUsername        = "ticket_username"
)

var (
	ticketRelatedTables = []string{
		ticketExtraAttributes,
		ticketGrantTicket,
		ticketIOU,
		ticketLastReferenced,
		ticketService,
		ticketType,
		ticketUsername,
	}
)

// delete queries block
var (
	deleteAllRelatedTables ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		for i := range ticketRelatedTables {
			err := func() error {
				log.Printf("delete from %s", ticketRelatedTables[i])
				stmt, err := tx.Prepare(fmt.Sprintf(deleteSomeRelatedTable, ticketRelatedTables[i]))
				if err != nil {
					return err
				}
				defer stmt.Close()
				_, err = stmt.Exec(t.ID)
				return err
			}()
			if err != nil {
				return err
			}
		}
		return nil
	}

	//  ticket delete by this
	deleteTicket ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		stmt, err := tx.Prepare(deleteTicketQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID)
		return err
	}
)

package ticket

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/deadcheat/cashew"
)

// deleters
var deleteServiceAccessors = []ticketAccessor{
	deleteAllRelatedTables,
	deleteGrantedServiceTicket,
	deleteTicket,
}

// deleters
var deleteGrantingTicketAccessors = []ticketAccessor{
	deleteAllRelatedTables,
	deleteGrantingTicket,
	deleteTicket,
}

const (
	ticketExtraAttributes = "ticket_extra_attribute"
	ticketIOU             = "ticket_iou"
	ticketLastReferenced  = "ticket_last_referenced"
	ticketService         = "ticket_service"
	ticketType            = "ticket_type"
	ticketPrimary         = "ticket_primary"
	ticketUsername        = "ticket_username"
)

var (
	ticketRelatedTables = []string{
		ticketExtraAttributes,
		ticketIOU,
		ticketLastReferenced,
		ticketService,
		ticketType,
		ticketPrimary,
		ticketUsername,
	}
)

// delete queries block
var (
	deleteAllRelatedTables ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		for i := range ticketRelatedTables {
			err := func() error {
				log.Printf("delete from %s", ticketRelatedTables[i])
				stmt, err := tx.Prepare(fmt.Sprintf(deleteSomeRelatedTableQeury, ticketRelatedTables[i]))
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
	// delete  granted service ticeket
	deleteGrantedServiceTicket ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		stmt, err := tx.Prepare(deleteGrantedServiceTicketQeury)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID)
		return err
	}

	// delete  granting ticeket relation
	deleteGrantingTicket ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		stmt, err := tx.Prepare(deleteGrantingTicketQeury)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID)
		return err
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

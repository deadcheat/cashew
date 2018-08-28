package ticket

import (
	"database/sql"

	"github.com/deadcheat/cashew"
	"github.com/deadcheat/cashew/timer"
	"github.com/deadcheat/cashew/values/errs"
	"github.com/rs/xid"
)

// inserter list
var insertAccessors = []ticketAccessor{
	insertTicket,
	insertTicketType,
	insertTicketService,
	insertTicketGrant,
	insertTicketUserName,
	insertTicketIOU,
	insertTicketLastReferenced,
	insertTicketExtraAttributes,
}

// insert queries block
var (
	// ticket all ticket insert by this
	insertTicket ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		stmt, err := tx.Prepare(createTicketQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.ClientHostName)
		return err
	}

	// inserter for ticket_type
	insertTicketType ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		stmt, err := tx.Prepare(createTicketTypeQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.Type.String())
		return err
	}

	// inserter for ticket_service
	insertTicketService ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		if t.Type != cashew.TicketTypeService {
			return nil
		}
		stmt, err := tx.Prepare(createTicketServiceQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.Service)
		return err
	}

	// inserter for ticket_username
	insertTicketUserName ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		if t.Type != cashew.TicketTypeService && t.Type != cashew.TicketTypeTicketGranting {
			return nil
		}
		stmt, err := tx.Prepare(createTicketUsernameQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.UserName)
		return err
	}

	// inserter for ticket_iou
	insertTicketIOU ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		if t.Type != cashew.TicketTypeProxyGranting {
			return nil
		}
		stmt, err := tx.Prepare(createTicketIOUQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.IOU)
		return err
	}
	// inserter for ticket_expires
	insertTicketLastReferenced ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
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
		stmt, err := tx.Prepare(createTicketLastReferencedQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.LastReferencedAt)
		return err
	}

	// inserter for ticket_extra_attributes
	insertTicketExtraAttributes ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		if t.Type != cashew.TicketTypeTicketGranting {
			return nil
		}
		stmt, err := tx.Prepare(createTicketExtraAttributesQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.ExtraAttributes)
		return err
	}

	// inserter for ticket_granting_ticket
	insertTicketGrant ticketAccessor = func(tx *sql.Tx, t *cashew.Ticket) error {
		if t.Type != cashew.TicketTypeService {
			return nil
		}
		if t.GrantedBy == nil {
			return errs.ErrTicketGrantedTicketIsNotFound
		}
		stmt, err := tx.Prepare(createTicketGrantQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(xid.New().String(), t.GrantedBy.ID, t.ID)
		return err
	}
)

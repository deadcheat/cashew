package ticket

import (
	"database/sql"
	"fmt"

	"github.com/deadcheat/cashew/values/errs"

	"github.com/deadcheat/cashew"
)

// Repository hold db connection and statements
type Repository struct {
	db *sql.DB
}

// New create new TicketRepository
func New(db *sql.DB) cashew.TicketRepository {
	return &Repository{db: db}
}

type ticketAccessor func(tx *sql.Tx, t *cashew.Ticket) error

// Create create new ticket
func (r *Repository) Create(t *cashew.Ticket) error {
	return r.executeOnNewTx(insertAccessors, t)
}

// Delete from all ticket-related table and ticket table
func (r *Repository) Delete(t *cashew.Ticket) error {
	fmt.Println(t.ID, t.Type)
	switch t.Type {
	case cashew.TicketTypeTicketGranting, cashew.TicketTypeProxyGranting:
		return r.executeOnNewTx(deleteGrantingTicketAccessors, t)
	case cashew.TicketTypeLogin, cashew.TicketTypeService:
		return r.executeOnNewTx(deleteServiceAccessors, t)
	}
	return errs.ErrInvalidTicketType
}

func (r *Repository) executeOnNewTx(accessors []ticketAccessor, t *cashew.Ticket) (err error) {
	var tx *sql.Tx
	tx, err = r.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	if err = r.executeTicketAccessors(tx, accessors, t); err != nil {
		return
	}
	return tx.Commit()
}

func (r *Repository) executeWithTx(tx *sql.Tx, accessors []ticketAccessor, t *cashew.Ticket) (err error) {
	return r.executeTicketAccessors(tx, accessors, t)
}

func (r *Repository) executeTicketAccessors(tx *sql.Tx, accessors []ticketAccessor, t *cashew.Ticket) (err error) {
	// FIXME if executer process increased wait-queues
	for i := range accessors {
		if err = accessors[i](tx, t); err != nil {
			return err
		}
	}
	return nil
}

// Find search for new ticket by ticket id
func (r *Repository) Find(id string) (t *cashew.Ticket, err error) {
	var gid string
	t, gid, err = r.findTicket(id)

	previous := t
	var g *cashew.Ticket
	for gid != "" {
		g, gid, err = r.findTicket(gid)
		if err != nil {
			break
		}
		previous.GrantedBy = g
		previous = g
	}
	return
}

func (r *Repository) findTicket(id string) (ticket *cashew.Ticket, granterTicketID string, err error) {
	var stmt *sql.Stmt
	stmt, err = r.db.Prepare(selectByTicketIDQuery)
	if err != nil {
		return nil, "", err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	ticket = new(cashew.Ticket)
	var (
		typeStr         sql.NullString
		service         sql.NullString
		grantedBy       sql.NullString
		userName        sql.NullString
		iou             sql.NullString
		extraAttributes interface{}
		primaryTicket   sql.NullString
	)
	err = row.Scan(
		&ticket.ID,
		&typeStr,
		&ticket.ClientHostName,
		&ticket.CreatedAt,
		&ticket.LastReferencedAt,
		&service,
		&userName,
		&iou,
		&extraAttributes,
		&grantedBy,
		&primaryTicket,
	)
	if err != nil {
		return nil, "", err
	}

	if typeStr.Valid {
		// NullString always return nil as error
		tmp, _ := typeStr.Value()
		ticket.Type = cashew.ParseTicketType(tmp.(string))
	}

	if service.Valid {
		// NullString always return nil as error
		tmp, _ := service.Value()
		ticket.Service = tmp.(string)
	}

	if iou.Valid {
		// NullString always return nil as error
		tmp, _ := iou.Value()
		ticket.IOU = tmp.(string)
	}

	if userName.Valid {
		// NullString always return nil as error
		tmp, _ := userName.Value()
		ticket.UserName = tmp.(string)
	}

	if grantedBy.Valid {
		tmp, _ := grantedBy.Value()
		granterTicketID, _ = tmp.(string)
	}

	ticket.Primary = primaryTicket.Valid
	return
}

// findAllRelatedTicket search for new ticket by ticket id
func (r *Repository) findAllRelatedTicket(id string) (ts []*cashew.Ticket, err error) {
	var stmt *sql.Stmt
	stmt, err = r.db.Prepare(selectAllTicketRelatedByGrantTicket)
	if err != nil {
		return
	}
	defer stmt.Close()
	var rows *sql.Rows
	rows, err = stmt.Query(id)
	if err != nil {
		return
	}
	defer rows.Close()
	ts = make([]*cashew.Ticket, 0)
	for rows.Next() {
		var (
			ticket          cashew.Ticket
			typeStr         sql.NullString
			service         sql.NullString
			userName        sql.NullString
			iou             sql.NullString
			extraAttributes interface{}
			primaryTicket   sql.NullString
		)
		err = rows.Scan(
			&ticket.ID,
			&typeStr,
			&ticket.ClientHostName,
			&ticket.CreatedAt,
			&ticket.LastReferencedAt,
			&service,
			&userName,
			&iou,
			&extraAttributes,
			&primaryTicket,
		)
		if err != nil {
			return nil, err
		}

		if typeStr.Valid {
			// NullString always return nil as error
			tmp, _ := typeStr.Value()
			ticket.Type = cashew.ParseTicketType(tmp.(string))
		}

		if service.Valid {
			// NullString always return nil as error
			tmp, _ := service.Value()
			ticket.Service = tmp.(string)
		}

		if iou.Valid {
			// NullString always return nil as error
			tmp, _ := iou.Value()
			ticket.IOU = tmp.(string)
		}

		if userName.Valid {
			// NullString always return nil as error
			tmp, _ := userName.Value()
			ticket.UserName = tmp.(string)
		}

		ticket.Primary = primaryTicket.Valid
		ts = append(ts, &ticket)
	}
	return ts, nil
}

func (r *Repository) deleteServiceTicket(tx *sql.Tx, t *cashew.Ticket) error {
	return r.executeWithTx(tx, deleteServiceAccessors, t)
}

func (r *Repository) deleteGrantingTicket(tx *sql.Tx, t *cashew.Ticket) error {
	return r.executeWithTx(tx, deleteGrantingTicketAccessors, t)
}

// DeleteRelatedTicket search for new ticket by ticket id
func (r *Repository) DeleteRelatedTicket(t *cashew.Ticket) (err error) {
	ts := []*cashew.Ticket{t}
	index := 0
	sentinel := ts[index]
	for sentinel != nil {
		var children []*cashew.Ticket
		// find all tickets granting this ticket
		children, _ = r.findAllRelatedTicket(sentinel.ID)
		if len(children) > 0 {
			ts = append(ts, children...)
		}
		if len(ts) <= index {
			break
		}
		sentinel = ts[index]
		index++
	}
	// start tran
	var tx *sql.Tx
	tx, err = r.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	// delete all tickets
	for i := len(ts) - 1; i >= 0; i-- {
		ti := ts[i]
		err = r.deleteServiceTicket(tx, ti)
		if err != nil {
			return
		}
	}
	// finally, delete ticket granting ticket
	r.deleteGrantingTicket(tx, t)
	return tx.Commit()
}

// Consume update last_referenced_time for ticket
func (r *Repository) Consume(t *cashew.Ticket) (err error) {
	if t == nil {
		return errs.ErrInvalidMethodCall
	}
	// start tran
	var tx *sql.Tx
	tx, err = r.db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	var findStmt *sql.Stmt
	findStmt, err = tx.Prepare(selectConsumedQuery)
	if err != nil {
		return
	}
	defer findStmt.Close()
	row := findStmt.QueryRow(t.ID)
	var one int
	err = row.Scan(&one)
	accessor := updateTicketLastReferenced
	switch err {
	case nil:
	case sql.ErrNoRows:
		accessor = insertTicketLastReferenced
	default:
		return
	}
	// insert when any records had not been found
	err = accessor(tx, t)
	if err != nil {
		return
	}
	return tx.Commit()
}

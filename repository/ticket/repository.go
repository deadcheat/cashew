package ticket

import (
	"database/sql"

	"github.com/deadcheat/cashew"
)

// Repository hold db connection and statements
type Repository struct {
	db *sql.DB
}

// New create new TicketRepository
func New(db *sql.DB) cashew.TicketRepository {
	return &Repository{db}
}

type ticketAccessor func(tx *sql.Tx, t *cashew.Ticket) error

// Create create new ticket
func (r *Repository) Create(t *cashew.Ticket) error {
	return r.executeTicketAccessors(insertAccessors, t)
}

// Delete from all ticket-related table and ticket table
func (r *Repository) Delete(t *cashew.Ticket) error {
	return r.executeTicketAccessors(deleteAccessors, t)
}

func (r *Repository) executeTicketAccessors(accessors []ticketAccessor, t *cashew.Ticket) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// FIXME if executer process increased wait-queues
	for i := range accessors {
		if err = accessors[i](tx, t); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Find search for new ticket by ticket id
func (r *Repository) Find(id string) (*cashew.Ticket, error) {

	stmt, err := r.db.Prepare(selectByTicketIDQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	var (
		ticket          cashew.Ticket
		typeStr         sql.NullString
		service         sql.NullString
		grantedBy       sql.NullString
		userName        sql.NullString
		iou             sql.NullString
		extraAttributes interface{}
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

	// TODO confirm this recursive call will not cause any problems
	if grantedBy.Valid {
		tmp, _ := grantedBy.Value()
		grantedByID, _ := tmp.(string)
		ticket.GrantedBy, err = r.Find(grantedByID)
		if err != nil {
			return nil, err
		}
	}

	return &ticket, nil
}

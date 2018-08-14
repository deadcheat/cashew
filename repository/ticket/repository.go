package ticket

import (
	"database/sql"
	"errors"

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

var logics = []ticketInserter{
	ticket,
	ticketType,
	ticketService,
	ticketGrant,
	ticketUserName,
	ticketIOU,
	ticketExpires,
	ticketExtraAttributes,
}

// Create create new ticket
func (r *Repository) Create(t *cashew.Ticket) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, logic := range logics {
		if err = logic(tx, t); err != nil {
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
		&ticket.ExpiresAt,
		&service,
		&userName,
		&iou,
		extraAttributes,
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

	if grantedBy.Valid {
		tmp, _ := grantedBy.Value()
		grantedByID, _ := tmp.(string)
		ticket.GrantedBy, err = r.Find(grantedByID)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

type ticketInserter func(tx *sql.Tx, t *cashew.Ticket) error

var (
	// ErrTicketGrantedTicketIsNotFound error when ticket has no granter
	ErrTicketGrantedTicketIsNotFound = errors.New("granting ticket is not found")
)

var (
	// ticket all ticket insert by this
	ticket ticketInserter = func(tx *sql.Tx, t *cashew.Ticket) error {
		stmt, err := tx.Prepare(createTicketQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.ClientHostName)
		return err
	}

	// inserter for ticket_type
	ticketType ticketInserter = func(tx *sql.Tx, t *cashew.Ticket) error {
		stmt, err := tx.Prepare(createTicketTypeQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.Type.String())
		return err
	}

	// inserter for ticket_service
	ticketService ticketInserter = func(tx *sql.Tx, t *cashew.Ticket) error {
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
	ticketUserName ticketInserter = func(tx *sql.Tx, t *cashew.Ticket) error {
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
	ticketIOU ticketInserter = func(tx *sql.Tx, t *cashew.Ticket) error {
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
	ticketExpires ticketInserter = func(tx *sql.Tx, t *cashew.Ticket) error {
		if t.Type != cashew.TicketTypeService && t.Type != cashew.TicketTypeProxy && t.Type != cashew.TicketTypeProxyGranting {
			return nil
		}
		stmt, err := tx.Prepare(createTicketExpiresQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.ExpiresAt)
		return err
	}

	// inserter for ticket_extra_attributes
	ticketExtraAttributes ticketInserter = func(tx *sql.Tx, t *cashew.Ticket) error {
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
	ticketGrant ticketInserter = func(tx *sql.Tx, t *cashew.Ticket) error {
		if t.Type != cashew.TicketTypeProxyGranting && t.Type != cashew.TicketTypeTicketGranting {
			return nil
		}
		if t.GrantedBy == nil {
			return ErrTicketGrantedTicketIsNotFound
		}
		stmt, err := tx.Prepare(createTicketGrantQuery)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(t.ID, t.GrantedBy.ID, t.ID)
		return err
	}
)

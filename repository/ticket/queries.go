package ticket

const (
	createTicketQuery                = `INSERT INTO ticket_type (id, client_hostname, created_at) VALUES (?, ?, DEFAULT)`
	createTicketTypeQuery            = `INSERT INTO ticket_type (id, type, created_at) VALUES (?, ?, DEFAULT)`
	createTicketGrantQuery           = `INSERT INTO ticket_grant_ticket (id, source_ticket_id, destination_ticket_id, created_at) VALUES (?, ?, ?,DEFAULT)`
	createTicketServiceQuery         = `INSERT INTO ticket_service (id, service, created_at) VALUES (?, ?, DEFAULT)`
	createTicketUsernameQuery        = `INSERT INTO ticket_username (id, username, created_at) VALUES (?, ?, DEFAULT)`
	createTicketIOUQuery             = `INSERT INTO ticket_iou (id, iou, created_at) VALUES (?, ?, DEFAULT)`
	createTicketExpiresQuery         = `INSERT INTO ticket_expires (id, expires_at, created_at) VALUES (?, ?, DEFAULT)`
	createTicketExtraAttributesQuery = `INSERT INTO ticket_extra_attributes (id, extra_attribute, created_at) VALUES (?, ?, DEFAULT)`
)

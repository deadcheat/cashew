package ticket

const (
	// queries for insert
	createTicketQuery                = `INSERT INTO tickets (id, client_hostname, created_at) VALUES (?, ?, DEFAULT)`
	createTicketTypeQuery            = `INSERT INTO ticket_type (ticket_id, type, created_at) VALUES (?, ?, DEFAULT)`
	createTicketGrantQuery           = `INSERT INTO ticket_grant_ticket (source_ticket_id, destination_ticket_id, created_at) VALUES (?, ?,DEFAULT)`
	createTicketServiceQuery         = `INSERT INTO ticket_service (ticket_id, service, created_at) VALUES (?, ?, DEFAULT)`
	createTicketUsernameQuery        = `INSERT INTO ticket_username (ticket_id, username, created_at) VALUES (?, ?, DEFAULT)`
	createTicketIOUQuery             = `INSERT INTO ticket_iou (ticket_id, iou, created_at) VALUES (?, ?, DEFAULT)`
	createTicketLastReferencedQuery  = `INSERT INTO ticket_last_referenced (ticket_id, last_referenced_at, created_at) VALUES (?, ?, DEFAULT)`
	createTicketExtraAttributesQuery = `INSERT INTO ticket_extra_attribute (ticket_id, extra_attribute, created_at) VALUES (?, ?, DEFAULT)`

	// queries for delete
	deleteSomeRelatedTable           = `DELETE FROM %s WHERE ticket_id = ?`
	deleteTicketQuery                = `DELETE FROM tickets WHERE id = ?`
	deleteTicketTypeQuery            = `DELETE FROM ticket_type WHERE ticket_id = ?`
	deleteTicketGrantQuery           = `DELETE FROM ticket_grant_ticket WHERE ticket_id = ?`
	deleteTicketServiceQuery         = `DELETE FROM ticket_service WHERE ticket_id = ?`
	deleteTicketUsernameQuery        = `DELETE FROM ticket_username WHERE ticket_id = ?`
	deleteTicketIOUQuery             = `DELETE FROM ticket_iou WHERE ticket_id = ?`
	deleteTicketExpiresQuery         = `DELETE FROM ticket_expires WHERE ticket_id = ?`
	deleteTicketLastReferencedQuery  = `DELETE FROM ticket_last_referenced WHERE ticket_id = ?`
	deleteTicketExtraAttributesQuery = `DELETE FROM ticket_extra_attribute WHERE ticket_id = ?`

	// queries for select
	selectByTicketIDQuery = `SELECT
    t.id,
    tt.type,
    t.client_hostname,
    t.created_at,
    tlr.last_referenced_at,
    ts.service,
    tu.username,
    i.iou,
    tea.extra_attribute,
    gt.id granted_by
  FROM tickets T
    LEFT JOIN ticket_type tt ON t.id = tt.ticket_id
    LEFT JOIN ticket_service ts ON t.id = ts.ticket_id
    LEFT JOIN ticket_last_referenced tlr ON t.id = tlr.ticket_id
    LEFT JOIN ticket_iou i ON t.id = i.ticket_id
    LEFT JOIN ticket_username tu ON t.id = tu.ticket_id
    LEFT JOIN ticket_extra_attribute tea ON t.id = tea.ticket_id
    LEFT JOIN ticket_grant_ticket tgt ON t.id = tgt.destination_ticket_id
    LEFT JOIN tickets gt ON tgt.source_ticket_id = gt.id
  WHERE
    t.id = ?
;
`
)

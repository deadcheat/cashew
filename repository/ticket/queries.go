package ticket

const (
	// queries for insert
	createTicketQuery                = `INSERT INTO tickets (id, client_hostname, created_at) VALUES (?, ?, DEFAULT)`
	createTicketTypeQuery            = `INSERT INTO ticket_type (ticket_id, type, created_at) VALUES (?, ?, DEFAULT)`
	createTicketGrantQuery           = `INSERT INTO ticket_grant_ticket (id, source_ticket_id, destination_ticket_id, created_at) VALUES (?, ?, ?,DEFAULT)`
	createTicketServiceQuery         = `INSERT INTO ticket_service (ticket_id, service, created_at) VALUES (?, ?, DEFAULT)`
	createTicketPrimaryQuery         = `INSERT INTO ticket_primary (ticket_id, created_at) VALUES (?, DEFAULT)`
	createTicketUsernameQuery        = `INSERT INTO ticket_username (ticket_id, username, created_at) VALUES (?, ?, DEFAULT)`
	createTicketIOUQuery             = `INSERT INTO ticket_iou (ticket_id, iou, created_at) VALUES (?, ?, DEFAULT)`
	createTicketLastReferencedQuery  = `INSERT INTO ticket_last_referenced (ticket_id, last_referenced_at, created_at) VALUES (?, ?, DEFAULT)`
	createTicketExtraAttributesQuery = `INSERT INTO ticket_extra_attribute (ticket_id, extra_attribute, created_at) VALUES (?, ?, DEFAULT)`

	// queries for delete
	deleteSomeRelatedTableQeury     = `DELETE FROM %s WHERE ticket_id = ?`
	deleteGrantedServiceTicketQeury = `DELETE FROM ticket_grant_ticket WHERE destination_ticket_id = ?`
	deleteGrantingTicketQeury       = `DELETE FROM ticket_grant_ticket WHERE source_ticket_id = ?`
	deleteTicketQuery               = `DELETE FROM tickets WHERE id = ?`

	// queries for update
	updateConsumeQuery = `UPDATE ticket_last_referenced tlr SET last_referenced_at = ? WHERE tlr.ticket_id = ?`

	// query for select last_referenced
	selectConsumedQuery = `SELECT 1 FROM ticket_last_referenced tlr WHERE tlr.ticket_id = ?`

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
    gt.id granted_by,
    tp.ticket_id as has_primary
  FROM tickets t
    LEFT JOIN ticket_type tt ON t.id = tt.ticket_id
    LEFT JOIN ticket_service ts ON t.id = ts.ticket_id
    LEFT JOIN ticket_last_referenced tlr ON t.id = tlr.ticket_id
    LEFT JOIN ticket_iou i ON t.id = i.ticket_id
    LEFT JOIN ticket_username tu ON t.id = tu.ticket_id
    LEFT JOIN ticket_extra_attribute tea ON t.id = tea.ticket_id
    LEFT JOIN ticket_grant_ticket tgt ON t.id = tgt.destination_ticket_id
    LEFT JOIN tickets gt ON tgt.source_ticket_id = gt.id
    LEFT JOIN ticket_primary tp ON t.id = tp.ticket_id
  WHERE
    t.id = ?
;
`

	selectAllTicketRelatedByGrantTicket = `SELECT
    t.id,
    tt.type,
    t.client_hostname,
    t.created_at,
    tlr.last_referenced_at,
    ts.service,
    tu.username,
    i.iou,
    tea.extra_attribute,
    tp.ticket_id as has_primary
  FROM tickets t
    LEFT JOIN ticket_type tt ON t.id = tt.ticket_id
    LEFT JOIN ticket_service ts ON t.id = ts.ticket_id
    LEFT JOIN ticket_last_referenced tlr ON t.id = tlr.ticket_id
    LEFT JOIN ticket_iou i ON t.id = i.ticket_id
    LEFT JOIN ticket_username tu ON t.id = tu.ticket_id
    LEFT JOIN ticket_extra_attribute tea ON t.id = tea.ticket_id
    LEFT JOIN ticket_grant_ticket tgt ON t.id = tgt.destination_ticket_id
    LEFT JOIN ticket_primary tp ON t.id = tp.ticket_id
  WHERE
    tgt.source_ticket_id = ?
`
)

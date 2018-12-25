package expiration

var (
	selectByTimeoutTicketQuery = `SELECT
    t.id,
    tt.type,
    t.created_at, 
    tlr.last_referenced_at
  FROM tickets t 
    LEFT JOIN ticket_type tt ON t.id = tt.ticket_id 
    LEFT JOIN ticket_last_referenced tlr ON t.id = tlr.ticket_id 
    LEFT JOIN ticket_username tu ON t.id = tu.ticket_id 
    LEFT JOIN ticket_grant_ticket tgt ON t.id = tgt.destination_ticket_id 
    LEFT JOIN tickets gt ON tgt.source_ticket_id = gt.id 
  WHERE
    tt.type IN( %s ) 
  AND (
    tlr.last_referenced_at <= ? 
    OR
    t.created_at <= ?
  )
;
  `
)

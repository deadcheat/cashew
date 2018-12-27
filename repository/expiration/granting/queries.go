package granting

var (
	selectByTimeoutTicketQuery = `SELECT
    t.id,
    tt.type,
    t.created_at, 
    tlr.last_referenced_at
  FROM tickets t 
    LEFT JOIN ticket_type tt ON t.id = tt.ticket_id 
    LEFT JOIN ticket_last_referenced tlr ON t.id = tlr.ticket_id 
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

package service

var (
	selectByTimeoutTicketQuery = `SELECT
    t.id,
    tt.type,
    t.created_at
  FROM tickets t 
    LEFT JOIN ticket_type tt ON t.id = tt.ticket_id 
  WHERE
    tt.type IN( %s ) 
  AND
    t.created_at <= ?
;
  `
)

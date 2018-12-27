package login

var (
	selectByTimeoutTicketQuery = `SELECT
    t.id,
    tt.type,
    t.created_at 
  FROM tickets t 
    LEFT JOIN ticket_type tt ON t.id = tt.ticket_id 
  WHERE
    tt.type = ?
  AND
    t.created_at <= ?
;
  `
)

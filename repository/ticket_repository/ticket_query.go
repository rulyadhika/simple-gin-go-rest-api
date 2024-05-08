package ticketrepository

const createNewTicketQuery = `INSERT INTO tickets(title, description, priority, status, created_by) 
	VALUES($1,$2,$3,$4,$5) RETURNING id, created_at, updated_at`

const findAllTicketQuery = `SELECT tickets.id, title, description, priority, status, tickets.created_at, tickets.updated_at,
	a.username, a.email, b.username, b.email, c.username, c.email
	FROM tickets JOIN users AS a ON tickets.created_by = a.id
	LEFT JOIN users AS b ON tickets.assign_to = b.id
	LEFT JOIN users AS c ON tickets.assign_by = c.id`

const findAllTicketByUserIdQuery = `SELECT tickets.id, title, description, priority, status, tickets.created_at, tickets.updated_at,
	a.username, a.email, b.username, b.email, c.username, c.email
	FROM tickets JOIN users AS a ON tickets.created_by = a.id
	LEFT JOIN users AS b ON tickets.assign_to = b.id
	LEFT JOIN users AS c ON tickets.assign_by = c.id 
	WHERE tickets.created_by = $1`

const findOneTicketByTicketIdQuery = `SELECT tickets.id, title, description, priority, status, tickets.created_at, tickets.updated_at,
	a.username, a.email, b.username, b.email, c.username, c.email
	FROM tickets JOIN users AS a ON tickets.created_by = a.id
	LEFT JOIN users AS b ON tickets.assign_to = b.id
	LEFT JOIN users AS c ON tickets.assign_by = c.id 
	WHERE tickets.id = $1`

const assignTicketToUserQuery = `UPDATE tickets SET assign_to=$1, assign_by=$2, status=$3, updated_at=$4 WHERE id=$5 RETURNING id`

const updateTicketStatusQuery = `UPDATE tickets SET status=$1, updated_at=$2 WHERE id=$3 RETURNING id`

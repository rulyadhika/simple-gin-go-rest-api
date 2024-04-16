CREATE TYPE TicketPriority AS ENUM('low','med','high','critical');
CREATE TYPE TicketStatus AS ENUM('open','in progress','resolved','closed');

CREATE TABLE IF NOT EXISTS tickets(
    id SERIAL PRIMARY KEY NOT NULL,
    ticket_id VARCHAR (255) NOT NULL UNIQUE,
    title VARCHAR (255) NOT NULL,
    description TEXT NOT NULL,
    priority TicketPriority NOT NULL,
    status TicketStatus NOT NULL,
    assigned_to INT DEFAULT NULL,
    assigned_by  INT DEFAULT NULL,
    created_by   INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY(assigned_to)
        REFERENCES users(id),
    FOREIGN KEY(assigned_by)
        REFERENCES users(id),
    FOREIGN KEY(created_by)
        REFERENCES users(id)
);
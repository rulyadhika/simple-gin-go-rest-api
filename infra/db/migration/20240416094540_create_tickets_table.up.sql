CREATE TYPE TicketPriority AS ENUM('low','med','high','critical');
CREATE TYPE TicketStatus AS ENUM('open','in progress','resolved','closed');

CREATE TABLE IF NOT EXISTS tickets(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    ticket_id VARCHAR (255) NOT NULL UNIQUE,
    title VARCHAR (255) NOT NULL,
    description TEXT NOT NULL,
    priority TicketPriority NOT NULL,
    status TicketStatus NOT NULL,
    assign_to UUID DEFAULT NULL,
    assign_by  UUID DEFAULT NULL,
    created_by   UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    FOREIGN KEY(assign_to)
        REFERENCES users(id),
    FOREIGN KEY(assign_by)
        REFERENCES users(id),
    FOREIGN KEY(created_by)
        REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS roles(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    role_name VARCHAR (50) UNIQUE NOT NULL
);
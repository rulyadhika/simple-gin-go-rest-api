CREATE TABLE IF NOT EXISTS roles(
    id SERIAL PRIMARY KEY NOT NULL,
    role_name VARCHAR (50) UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS users_roles(
    id SERIAL PRIMARY KEY NOT NULL,
    user_id INT NOT NULL,
    role_id INT NOT NULL,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    FOREIGN KEY(user_id)
        REFERENCES users(id)
            ON DELETE CASCADE,
    FOREIGN KEY(role_id)
        REFERENCES roles(id)
);
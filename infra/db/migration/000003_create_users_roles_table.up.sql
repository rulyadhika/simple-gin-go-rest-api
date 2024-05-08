CREATE TABLE IF NOT EXISTS users_roles(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now(),
    FOREIGN KEY(user_id)
        REFERENCES users(id)
            ON DELETE CASCADE,
    FOREIGN KEY(role_id)
        REFERENCES roles(id)
);
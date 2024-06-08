CREATE TABLE
    IF NOT EXISTS accounts_password_resets (
        user_id UUID NOT NULL,
        token VARCHAR(8) UNIQUE NOT NULL,
        request_time TIMESTAMPTZ DEFAULT now (),
        expiration_time TIMESTAMPTZ NOT NULL,
        next_request_available_at TIMESTAMPTZ NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );
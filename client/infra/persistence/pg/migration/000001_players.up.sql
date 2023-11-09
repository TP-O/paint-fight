CREATE TABLE players (
    id uuid PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(320) NOT NULL,
    password VARCHAR(72) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT 1,
    email_verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    password_updated_at TIMESTAMP
);

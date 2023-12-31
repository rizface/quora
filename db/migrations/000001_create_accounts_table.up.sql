CREATE TABLE IF NOT EXISTS accounts (
    id UUID NOT NULL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    email_is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMP DEFAULT NULL
)
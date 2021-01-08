CREATE TABLE users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    encrypted_password VARCHAR NOT NULL
);

CREATE TABLE posts (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    title VARCHAR NOT NULL,
    content jsonb NOT NULL,
    owner BIGINT REFERENCES users(id)
);
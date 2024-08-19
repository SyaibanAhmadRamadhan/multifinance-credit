CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    password   VARCHAR(255),
    email      VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

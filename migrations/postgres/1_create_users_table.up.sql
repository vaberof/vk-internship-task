CREATE TABLE IF NOT EXISTS users
(
    id       SERIAL PRIMARY KEY,
    email    VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR            NOT NULL,
    role     VARCHAR(50)        NOT NULL CHECK ( role IN ('user', 'admin') )
);
CREATE INDEX IF NOT EXISTS email_idx ON users (email);

-- Password: asdf1234 --
INSERT INTO users(email, password, role)
VALUES ('user@example.com', '$2a$10$JT0HAAksN7kvv6m0TXAvIejUzNOs19uRA7Ae8qIjn5lLa2hP1isNK', 'user');

-- Password: asdf1234 --
INSERT INTO users(email, password, role)
VALUES ('admin@example.com', '$2a$10$JT0HAAksN7kvv6m0TXAvIejUzNOs19uRA7Ae8qIjn5lLa2hP1isNK', 'admin');




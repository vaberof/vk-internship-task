CREATE TABLE IF NOT EXISTS films
(
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(150) NOT NULL CHECK (LENGTH(title) > 0),
    description  VARCHAR(1000),
    release_date DATE         NOT NULL,
    rating       SMALLINT CHECK (rating BETWEEN 0 AND 10)
);
CREATE INDEX IF NOT EXISTS title_idx ON films (title);
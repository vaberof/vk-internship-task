CREATE TABLE IF NOT EXISTS actors
(
    id        SERIAL PRIMARY KEY,
    name      VARCHAR(100)                           NOT NULL,
    sex       SMALLINT CHECK ( sex IN (0, 1, 2, 9) ) NOT NULL DEFAULT 0,
    birthdate DATE                                   NOT NULL
);
CREATE INDEX IF NOT EXISTS name_idx ON actors (name);

COMMENT ON COLUMN actors.sex IS 'Codes for the representation of human sexes according to ISO/IEC 5218';

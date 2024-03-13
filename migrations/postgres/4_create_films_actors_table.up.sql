CREATE TABLE IF NOT EXISTS films_actors
(
    id       SERIAL PRIMARY KEY,
    film_id  INT NOT NULL REFERENCES films (id),
    actor_id INT NOT NULL REFERENCES actors (id)
);
CREATE INDEX IF NOT EXISTS film_id_idx ON films_actors (film_id);
CREATE INDEX IF NOT EXISTS actor_id_idx ON films_actors (actor_id);
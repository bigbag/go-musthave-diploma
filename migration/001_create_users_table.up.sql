CREATE TABLE IF NOT EXISTS users
(
    id              varchar primary key,
    password        varchar NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_idx ON users USING btree (id);
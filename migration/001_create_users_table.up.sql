CREATE TABLE IF NOT EXISTS users
(
    id              serial primary key,
    login           varchar NOT NULL,
    password        varchar NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_idx ON users USING btree (id);
CREATE UNIQUE INDEX IF NOT EXISTS users_login_uniq_idx ON users USING btree (login);
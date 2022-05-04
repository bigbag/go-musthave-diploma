CREATE TABLE IF NOT EXISTS orders
(
    id              varchar primary key,
    user_id         varchar NOT NULL,
    amount          decimal (15,2) DEFAULT 0 NOT NULL,
    CONSTRAINT      orders_amount_positive CHECK (amount >= 0),
    uploaded_at     timestamp,
    status          varchar DEFAULT 'NEW' NOT NULL,
    is_final        bool default false not null
);
CREATE UNIQUE INDEX IF NOT EXISTS orders_idx ON orders USING btree (id);
CREATE TABLE IF NOT EXISTS withdrawals
(
    id              varchar primary key,
    user_id         varchar NOT NULL,
    amount          decimal (15,2) DEFAULT 0 NOT NULL,
    CONSTRAINT      withdrawals_amount_positive CHECK (amount >= 0),
    processed_at    timestamp
);
CREATE UNIQUE INDEX IF NOT EXISTS withdrawals_idx ON withdrawals USING btree (id);

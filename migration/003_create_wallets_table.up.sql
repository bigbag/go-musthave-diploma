CREATE TABLE IF NOT EXISTS wallets
(
    id              serial primary key,
    user_id         varchar NOT NULL,
    balance         decimal (15,2) DEFAULT 0 NOT NULL,
    CONSTRAINT      wallets_balance_positive CHECK (balance >= 0),
    withdrawal      decimal (15,2) DEFAULT 0 NOT NULL,
    CONSTRAINT      wallets_withdrawal_positive CHECK (withdrawal >= 0)
);
CREATE UNIQUE INDEX IF NOT EXISTS wallets_idx ON wallets USING btree (id);
CREATE UNIQUE INDEX IF NOT EXISTS wallets_user_id_uniq_idx ON wallets USING btree (user_id);
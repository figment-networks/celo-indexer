CREATE TABLE IF NOT EXISTS account_activity_sequences
(
    id               BIGSERIAL                NOT NULL,

    height           DECIMAL(65, 0)           NOT NULL,
    time             TIMESTAMP WITH TIME ZONE NOT NULL,

    transaction_hash TEXT                     NOT NULL,
    address          TEXT                     NOT NULL,
    amount           DECIMAL(65, 0)           NOT NULL,
    kind             TEXT                     NOT NULL,
    data             JSONB                    NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_account_activity_sequences_height on account_activity_sequences (height);
CREATE index idx_account_activity_sequences_transaction_hash on account_activity_sequences (transaction_hash);
CREATE index idx_account_activity_sequences_address_kind on account_activity_sequences (address, kind);
CREATE index idx_account_activity_sequences_address on account_activity_sequences (address);
CREATE index idx_account_activity_sequences_kind on account_activity_sequences (kind);

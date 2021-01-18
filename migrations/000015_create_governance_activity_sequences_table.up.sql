CREATE TABLE IF NOT EXISTS governance_activity_sequences
(
    id               BIGSERIAL                NOT NULL,

    height           DECIMAL(65, 0)           NOT NULL,
    time             TIMESTAMP WITH TIME ZONE NOT NULL,

    transaction_hash TEXT                     NOT NULL,
    proposal_id      DECIMAL(65, 0)           NOT NULL,
    account          TEXT,
    kind             TEXT                     NOT NULL,
    data             JSONB                    NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_governance_activity_sequences_height on governance_activity_sequences (height);
CREATE index idx_governance_activity_sequences_transaction_hash on governance_activity_sequences (transaction_hash);
CREATE index idx_governance_activity_sequences_proposal_id on governance_activity_sequences (proposal_id);
CREATE index idx_governance_activity_sequences_kind on governance_activity_sequences (kind);

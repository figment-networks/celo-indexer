CREATE TABLE IF NOT EXISTS validator_group_sequences
(
    id                BIGSERIAL                NOT NULL,

    height            DECIMAL(65, 0)           NOT NULL,
    time              TIMESTAMP WITH TIME ZONE NOT NULL,

    address           TEXT                     NOT NULL,
    commission        DECIMAL(65, 0)           NOT NULL,
    active_votes      DECIMAL(65, 0)           NOT NULL,
    active_vote_units DECIMAL(65, 0)           NOT NULL,
    pending_votes     DECIMAL(65, 0)           NOT NULL,
    members_count     INT                      NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_validator_group_sequences_height on validator_group_sequences (height);
CREATE index idx_validator_group_sequences_time on validator_group_sequences (time);
CREATE index idx_validator_group_sequences_address on validator_group_sequences (address);

CREATE TABLE IF NOT EXISTS validator_sequences
(
    id           BIGSERIAL                NOT NULL,

    height       DECIMAL(65, 0)           NOT NULL,
    time         TIMESTAMP WITH TIME ZONE NOT NULL,

    address      TEXT                     NOT NULL,
    affiliation  TEXT                     NOT NULL,
    signed       BOOLEAN,
    score        DECIMAL(65, 0)           NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_validator_sequences_height on validator_sequences (height);
CREATE index idx_validator_sequences_time on validator_sequences (time);
CREATE index idx_validator_sequences_address on validator_sequences (address);

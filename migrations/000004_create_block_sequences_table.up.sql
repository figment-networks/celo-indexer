CREATE TABLE IF NOT EXISTS block_sequences
(
    id               BIGSERIAL                NOT NULL,

    height           DECIMAL(65, 0)           NOT NULL,
    time             TIMESTAMP WITH TIME ZONE NOT NULL,

    tx_count         DOUBLE PRECISION,
    size             DOUBLE PRECISION,
    gas_used         DECIMAL(65, 0),
    total_difficulty DECIMAL(65, 0),

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_block_sequences_height on block_sequences (height);
CREATE index idx_block_sequences_time on block_sequences (time);

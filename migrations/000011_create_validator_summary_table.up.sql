CREATE TABLE IF NOT EXISTS validator_summary
(
    id            BIGSERIAL                NOT NULL,

    time_interval VARCHAR                  NOT NULL,
    time_bucket   TIMESTAMP WITH TIME ZONE NOT NULL,
    index_version INT                      NOT NULL,

    address       TEXT                     NOT NULL,
    score_avg     DECIMAL(65, 0)           NOT NULL,
    score_max     DECIMAL(65, 0)           NOT NULL,
    score_min     DECIMAL(65, 0)           NOT NULL,

    signed_avg    DECIMAL                  NOT NULL,
    signed_min    BIGINT                   NOT NULL,
    signed_max    BIGINT                   NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_validator_summary_time on validator_summary (time_interval, time_bucket);
CREATE index idx_validator_summary_index_version on validator_summary (index_version);
CREATE index idx_validator_summary_address on validator_summary (address);
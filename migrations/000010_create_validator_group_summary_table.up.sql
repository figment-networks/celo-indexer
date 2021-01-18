CREATE TABLE IF NOT EXISTS validator_group_summary
(
    id                    BIGSERIAL                NOT NULL,

    time_interval         VARCHAR                  NOT NULL,
    time_bucket           TIMESTAMP WITH TIME ZONE NOT NULL,
    index_version         INT                      NOT NULL,

    address               TEXT                     NOT NULL,
    commission_avg        DECIMAL(65, 0)           NOT NULL,
    commission_max        DECIMAL(65, 0)           NOT NULL,
    commission_min        DECIMAL(65, 0)           NOT NULL,
    active_votes_avg      DECIMAL(65, 0)           NOT NULL,
    active_votes_max      DECIMAL(65, 0)           NOT NULL,
    active_votes_min      DECIMAL(65, 0)           NOT NULL,
    pending_votes_avg     DECIMAL(65, 0)           NOT NULL,
    pending_votes_max     DECIMAL(65, 0)           NOT NULL,
    pending_votes_min     DECIMAL(65, 0)           NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_validator_group_summary_time on validator_group_summary (time_interval, time_bucket);
CREATE index idx_validator_group_summary_index_version on validator_group_summary (index_version);
CREATE index idx_validator_group_summary_address on validator_group_summary (address);
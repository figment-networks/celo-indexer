CREATE TABLE IF NOT EXISTS validator_group_aggregates
(
    id                         BIGSERIAL                NOT NULL,
    created_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at                 TIMESTAMP WITH TIME ZONE NOT NULL,

    started_at_height          DECIMAL(65, 0)           NOT NULL,
    started_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    recent_at_height           DECIMAL(65, 0)           NOT NULL,
    recent_at                  TIMESTAMP WITH TIME ZONE NOT NULL,

    address                    TEXT                     NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_validator_group_aggregates_address on validator_group_aggregates (address);
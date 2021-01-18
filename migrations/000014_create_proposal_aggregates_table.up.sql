CREATE TABLE IF NOT EXISTS proposal_aggregates
(
    id                         BIGSERIAL                NOT NULL,
    created_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at                 TIMESTAMP WITH TIME ZONE NOT NULL,

    started_at_height          DECIMAL(65, 0)           NOT NULL,
    started_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    recent_at_height           DECIMAL(65, 0)           NOT NULL,
    recent_at                  TIMESTAMP WITH TIME ZONE NOT NULL,

    proposal_id                BIGINT                   NOT NULL,
    proposer_address           TEXT                     NOT NULL,
    description_url            TEXT,
    deposit                    TEXT                     NOT NULL,
    transaction_count          BIGINT                   NOT NULL,
    proposed_at                TIMESTAMP WITH TIME ZONE NOT NULL,
    proposed_at_height         DECIMAL(65, 0)           NOT NULL,

    recent_stage               TEXT,

    dequeue_address            TEXT,
    dequeued_at                TIMESTAMP WITH TIME ZONE,
    dequeued_at_height         DECIMAL(65, 0),

    approval_address           TEXT,
    approved_at                TIMESTAMP WITH TIME ZONE,
    approved_at_height         DECIMAL(65, 0),

    executor_address           TEXT,
    executed_at                TIMESTAMP WITH TIME ZONE,
    executed_at_height         DECIMAL(65, 0),

    expired_at                 TIMESTAMP WITH TIME ZONE,
    expired_at_height          DECIMAL(65, 0),

    upvotes_total              TEXT   DEFAULT '0',
    yes_votes_total            BIGINT DEFAULT 0,
    yes_votes_weight_total     TEXT   DEFAULT '0',
    no_votes_total             BIGINT DEFAULT 0,
    no_votes_weight_total      TEXT   DEFAULT '0',
    abstain_votes_total        BIGINT DEFAULT 0,
    abstain_votes_weight_total TEXT   DEFAULT '0',
    votes_total                BIGINT DEFAULT 0,
    votes_weight_total         TEXT   DEFAULT '0',

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_proposal_aggregates_proposal_id on proposal_aggregates (proposal_id);
CREATE index idx_proposal_aggregates_proposer_address on proposal_aggregates (proposer_address);
CREATE TABLE IF NOT EXISTS syncables
(
    id            BIGSERIAL                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

    chain_id      BIGINT                   NOT NULL,
    height        DECIMAL(65, 0)           NOT NULL,
    time          TIMESTAMP WITH TIME ZONE NOT NULL,
    epoch         DECIMAL(65, 0),
    last_in_epoch BOOLEAN,

    index_version INT                      NOT NULL,
    status        SMALLINT DEFAULT 0,
    report_id     BIGINT,
    started_at    TIMESTAMP WITH TIME ZONE,
    processed_at  TIMESTAMP WITH TIME ZONE,
    duration      BIGINT,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_syncables_report_id on syncables (report_id);
CREATE index idx_syncables_height on syncables (height);
CREATE index idx_syncables_epoch on syncables (epoch);
CREATE index idx_syncables_index_version on syncables (index_version);
CREATE index idx_syncables_processed_at on syncables (processed_at);

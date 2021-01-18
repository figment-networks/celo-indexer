CREATE UNIQUE INDEX idx_validator_sequences_height_address ON validator_sequences(height, address);
CREATE UNIQUE INDEX idx_validator_group_sequences_height_address ON validator_group_sequences(height, address);
CREATE UNIQUE INDEX idx_system_events_height_actor_kind ON system_events(height, actor, kind);

CREATE UNIQUE INDEX idx_validator_summary_multi ON validator_summary(time_interval, time_bucket, index_version, address);
CREATE UNIQUE INDEX idx_validator_group_summary_multi ON validator_group_summary(time_interval, time_bucket, index_version, address);
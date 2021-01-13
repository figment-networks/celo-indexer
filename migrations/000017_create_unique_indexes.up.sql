CREATE UNIQUE INDEX idx_validator_sequences_height_address ON validator_sequences(height, address);

CREATE UNIQUE INDEX idx_validator_group_sequences_height_address ON validator_group_sequences(height, address);
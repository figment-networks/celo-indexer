package model

type ValidatorAgg struct {
	*Model
	*Aggregate

	Address                 string `json:"address"`
	RecentName              string `json:"recent_name"`
	RecentMetadataUrl       string `json:"recent_metadata_url"`
	RecentAsValidatorHeight int64  `json:"recent_as_validator_height"`
	AccumulatedUptime       int64  `json:"accumulated_uptime"`
	AccumulatedUptimeCount  int64  `json:"accumulated_uptime_count"`
}

// - Methods
func (ValidatorAgg) TableName() string {
	return "validator_aggregates"
}

func (s *ValidatorAgg) Valid() bool {
	return s.Aggregate.Valid() &&
		s.Address != ""
}

func (s *ValidatorAgg) Equal(m ValidatorAgg) bool {
	return s.Address == m.Address
}

func (s *ValidatorAgg) Update(u *ValidatorAgg) {
	s.Aggregate.RecentAtHeight = u.Aggregate.RecentAtHeight
	s.Aggregate.RecentAt = u.Aggregate.RecentAt
	s.RecentName = u.RecentName
	s.RecentMetadataUrl = u.RecentMetadataUrl
	s.RecentAsValidatorHeight = u.RecentAsValidatorHeight
	s.AccumulatedUptime = u.AccumulatedUptime
	s.AccumulatedUptimeCount = u.AccumulatedUptimeCount
}

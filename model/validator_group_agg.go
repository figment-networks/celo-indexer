package model

type ValidatorGroupAgg struct {
	*ModelWithTimestamps
	*Aggregate

	Address                string `json:"address"`
	RecentName             string `json:"recent_name"`
	RecentMetadataUrl      string `json:"recent_metadata_url"`
	AccumulatedUptime      int64  `json:"accumulated_uptime"`
	AccumulatedUptimeCount int64  `json:"accumulated_uptime_count"`
}

func (ValidatorGroupAgg) TableName() string {
	return "validator_group_aggregates"
}

func (s *ValidatorGroupAgg) Valid() bool {
	return s.Aggregate.Valid() &&
		s.Address != ""
}

func (s *ValidatorGroupAgg) Equal(m ValidatorGroupAgg) bool {
	return s.Address == m.Address
}

func (s *ValidatorGroupAgg) Update(u *ValidatorGroupAgg) {
	s.Aggregate.RecentAtHeight = u.Aggregate.RecentAtHeight
	s.Aggregate.RecentAt = u.Aggregate.RecentAt
	s.RecentName = u.RecentName
	s.RecentMetadataUrl = u.RecentMetadataUrl
	s.AccumulatedUptimeCount = s.AccumulatedUptimeCount + u.AccumulatedUptimeCount
	s.AccumulatedUptime = s.AccumulatedUptime + u.AccumulatedUptime
}

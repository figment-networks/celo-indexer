package model

type ValidatorGroupAgg struct {
	*Model
	*Aggregate

	Address string `json:"address"`
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
}

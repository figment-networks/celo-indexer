package model

import "github.com/figment-networks/celo-indexer/types"

const (
	SystemEventGroupRewardChange1 SystemEventKind = "group_reward_change_1"
	SystemEventGroupRewardChange2 SystemEventKind = "group_reward_change_2"
	SystemEventGroupRewardChange3 SystemEventKind = "group_reward_change_3"
	SystemEventJoinedActiveSet    SystemEventKind = "joined_active_set"
	SystemEventLeftActiveSet      SystemEventKind = "left_active_set"
	SystemEventMissedNConsecutive SystemEventKind = "missed_n_consecutive"
	SystemEventMissedNofM         SystemEventKind = "missed_n_of_m"
)

type SystemEventKind string

func (o SystemEventKind) String() string {
	return string(o)
}

type SystemEvent struct {
	*Model

	Height int64           `json:"height"`
	Time   types.Time      `json:"time"`
	Actor  string          `json:"actor"`
	Kind   SystemEventKind `json:"kind"`
	Data   types.Jsonb     `json:"data"`
}

func (o SystemEvent) Update(m SystemEvent) {
	o.Height = m.Height
	o.Time = m.Time
	o.Actor = m.Actor
	o.Kind = m.Kind
	o.Data = m.Data
}

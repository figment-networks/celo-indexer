package model

type BlockSummary struct {
	*ModelWithTimestamps
	*Summary

	Count                      int64   `json:"count"`
	BlockTimeAvg               float64 `json:"block_time_avg"`
}

func (BlockSummary) TableName() string {
	return "block_summary"
}

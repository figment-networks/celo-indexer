package model

import "github.com/figment-networks/celo-indexer/types"

type ValidatorSummary struct {
	*Model
	*Summary

	Address string `json:"address"`

	ScoreAvg types.Quantity `json:"score_avg"`
	ScoreMin types.Quantity `json:"score_min"`
	ScoreMax types.Quantity `json:"score_max"`

	SignedAvg float64 `json:"signed_avg"`
	SignedMax int64   `json:"signed_max"`
	SignedMin int64   `json:"signed_min"`
}

func (ValidatorSummary) TableName() string {
	return "validator_summary"
}

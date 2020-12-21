package validator

import (
	"github.com/figment-networks/celo-indexer/model"
)

type AggListView struct {
	Items []model.ValidatorAgg `json:"items"`
}

func ToAggListView(ms []model.ValidatorAgg) *AggListView {
	return &AggListView{
		Items: ms,
	}
}

type AggDetailsView struct {
	*model.ModelWithTimestamps
	*model.Aggregate

	Address                 string  `json:"address"`
	RecentAsValidatorHeight int64   `json:"recent_as_validator_height"`
	RecentName              string  `json:"recent_name"`
	RecentMetadataUrl       string  `json:"recent_metadata_url"`
	Uptime                  float64 `json:"uptime"`

	LastSequences []model.ValidatorSeq `json:"last_sequences"`
}

func ToAggDetailsView(m *model.ValidatorAgg, validatorSequences []model.ValidatorSeq) *AggDetailsView {
	return &AggDetailsView{
		ModelWithTimestamps: m.ModelWithTimestamps,
		Aggregate:           m.Aggregate,

		Address:                 m.Address,
		RecentAsValidatorHeight: m.RecentAsValidatorHeight,
		RecentName:              m.RecentName,
		RecentMetadataUrl:       m.RecentMetadataUrl,
		Uptime:                  float64(m.AccumulatedUptime) / float64(m.AccumulatedUptimeCount),

		LastSequences: validatorSequences,
	}
}

type SeqListItem struct {
	*model.Sequence

	Address     string  `json:"address"`
	Name        string  `json:"name"`
	MetadataUrl string  `json:"metadata_url"`
	Signed      *bool   `json:"signed"`
	Score       float64 `json:"score"`
	Affiliation string  `json:"affiliation"`
}

type SeqListView struct {
	Items []SeqListItem `json:"items"`
}

func ToSeqListView(validatorSeqs []model.ValidatorSeq) SeqListView {
	var items []SeqListItem
	for _, m := range validatorSeqs {
		item := SeqListItem{
			Sequence: m.Sequence,

			Address:     m.Address,
			Name:        m.Name,
			MetadataUrl: m.MetadataUrl,
			Signed:      m.Signed,
			Score:       m.ScoreAsPercentage(),
			Affiliation: m.Affiliation,
		}

		items = append(items, item)
	}

	return SeqListView{
		Items: items,
	}
}

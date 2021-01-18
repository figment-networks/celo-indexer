package theceloclient

import "github.com/figment-networks/celo-indexer/types"

type Proposals struct {
	Items ProposalItems `json:"items"`
}

type ProposalItems map[string]ProposalDetails

type ProposalDetails struct {
	Status         string   `json:"status"`
	Timespan       int64    `json:"timespan"`
	Title          string   `json:"title"`
	DescriptionUrl string   `json:"descriptionUrl"`
	Proposer       Proposer `json:"proposer"`
	Upvoted        Upvoted  `json:"upvoted"`
	Dequeue        Dequeue  `json:"dequeue"`
	Approval       Approval `json:"approval"`
	Voted          Voted    `json:"voted"`
	Executed       Executed `json:"executed"`
}

type Proposer struct {
	Address   string         `json:"address"`
	Deposit   types.Quantity `json:"deposit"`
	Timestamp int64          `json:"timestamp"`
}

type Upvoted struct {
	Peoples int64 `json:"peoples"`
	Upvotes int64 `json:"upvotes"`
}

type Dequeue struct {
	Address   string `json:"address"`
	Timestamp int64  `json:"timestamp"`
}

type Approval struct {
	Address   string `json:"address"`
	Timestamp string  `json:"timestamp"`
}

type Voted struct {
	Peoples int64          `json:"peoples"`
	Weight  types.Quantity `json:"weight"`
}

type Executed struct {
	From            string `json:"from"`
	Timestamp       string `json:"timestamp"`
	BlockNumber     string `json:"block_number"`
	TransactionHash string `json:"transaction_hash"`
}

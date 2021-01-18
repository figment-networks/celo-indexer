package model

import "github.com/figment-networks/celo-indexer/types"

const (
	VoteAbstain = 1
	VoteNo      = 2
	VoteYes     = 3

	ProposalStageProposed = "proposed"
	ProposalStageApproved = "approved"
	ProposalStageExecuted = "executed"
	ProposalStageDequeued = "dequeued"
	ProposalStageExpired  = "expired"
)

type ProposalAgg struct {
	*ModelWithTimestamps
	*Aggregate

	ProposalId       uint64     `json:"proposal_id"`
	ProposerAddress  string     `json:"proposer_address"`
	DescriptionUrl   string     `json:"description_url"`
	Deposit          string     `json:"deposit"`
	TransactionCount int64      `json:"transaction_count"`
	ProposedAtHeight int64      `json:"proposed_at_height"`
	ProposedAt       types.Time `json:"proposed_at"`

	RecentStage string `json:"recent_stage"`

	DequeueAddress   string     `json:"dequeue_address"`
	DequeuedAtHeight int64      `json:"dequeued_at_height"`
	DequeuedAt       types.Time `json:"dequeued_at"`

	ApprovalAddress  string     `json:"approval_address"`
	ApprovedAtHeight int64      `json:"approved_at_height"`
	ApprovedAt       types.Time `json:"approved_at"`

	ExecutorAddress  string     `json:"executor_address"`
	ExecutedAtHeight int64      `json:"executed_at_height"`
	ExecutedAt       types.Time `json:"executed_at"`

	ExpiredAtHeight int64      `json:"expired_at_height"`
	ExpiredAt       types.Time `json:"expired_at"`

	UpvotesTotal            string `json:"upvotes_total"`
	YesVotesTotal           uint64 `json:"yes_votes_total"`
	YesVotesWeightTotal     string `json:"yes_votes_weight_total"`
	NoVotesTotal            uint64 `json:"no_votes_total"`
	NoVotesWeightTotal      string `json:"no_votes_weight_total"`
	AbstainVotesTotal       uint64 `json:"abstain_votes_total"`
	AbstainVotesWeightTotal string `json:"abstain_votes_weight_total"`
	VotesTotal              uint64 `json:"votes_total"`
	VotesWeightTotal        string `json:"votes_weight_total"`
}

func (ProposalAgg) TableName() string {
	return "proposal_aggregates"
}

func (s *ProposalAgg) Valid() bool {
	return s.Aggregate.Valid() &&
		s.ProposalId != 0
}

func (s *ProposalAgg) Equal(m ProposalAgg) bool {
	return s.ProposalId == m.ProposalId
}

func (s *ProposalAgg) Update(u *ProposalAgg) {
	s.Aggregate.RecentAtHeight = u.Aggregate.RecentAtHeight
	s.Aggregate.RecentAt = u.Aggregate.RecentAt

	s.ProposerAddress = u.ProposerAddress
	s.Deposit = u.Deposit
	s.TransactionCount = u.TransactionCount
	s.ProposedAtHeight = u.ProposedAtHeight
	s.ProposedAt = u.ProposedAt
	s.RecentStage = u.RecentStage
	s.DequeueAddress = u.DequeueAddress
	s.DequeuedAtHeight = u.DequeuedAtHeight
	s.DequeuedAt = u.DequeuedAt
	s.ApprovalAddress = u.ApprovalAddress
	s.ApprovedAtHeight = u.ApprovedAtHeight
	s.ApprovedAt = u.ApprovedAt
	s.ExecutorAddress = u.ExecutorAddress
	s.ExecutedAtHeight = u.ExecutedAtHeight
	s.ExecutedAt = u.ExecutedAt
	s.ExpiredAtHeight = u.ExpiredAtHeight
	s.ExpiredAt = u.ExpiredAt
	s.UpvotesTotal = u.UpvotesTotal
	s.YesVotesTotal = u.YesVotesTotal
	s.YesVotesWeightTotal = u.YesVotesWeightTotal
	s.NoVotesTotal = u.NoVotesTotal
	s.NoVotesWeightTotal = u.NoVotesWeightTotal
	s.AbstainVotesTotal = u.AbstainVotesTotal
	s.AbstainVotesWeightTotal = u.AbstainVotesWeightTotal
	s.VotesTotal = u.VotesTotal
	s.VotesWeightTotal = u.VotesWeightTotal
}

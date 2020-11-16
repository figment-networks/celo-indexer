package figmentclient

import (
	blscrypto "github.com/ethereum/go-ethereum/crypto/bls"
	"math/big"
)

const (
	OperationTypeInternalTransfer                 = "InternalTransfer"
	OperationTypeValidatorGroupVoteCast           = "ValidatorGroupVoteCast"
	OperationTypeValidatorGroupVoteActivated      = "ValidatorGroupVoteActivated"
	OperationTypeValidatorGroupPendingVoteRevoked = "ValidatorGroupPendingVoteRevoked"
	OperationTypeValidatorGroupActiveVoteRevoked  = "ValidatorGroupActiveVoteRevoked"
	OperationTypeAccountCreated                   = "AccountCreated"
	OperationTypeAccountSlashed                   = "AccountSlashed"
	OperationTypeVoteSignerAuthorized             = "VoteSignerAuthorized"
	OperationTypeValidatorSignerAuthorized        = "ValidatorSignerAuthorized"
	OperationTypeAttestationSignerAuthorized      = "AttestationSignerAuthorized"
	OperationTypeGoldLocked                       = "GoldLocked"
	OperationTypeGoldRelocked                     = "GoldRelocked"
	OperationTypeGoldUnlocked                     = "GoldUnlocked"
	OperationTypeGoldWithdrawn                    = "GoldWithdrawn"
	OperationTypeValidatorEpochPaymentDistributed = "ValidatorEpochPaymentDistributed"
)

type ChainStatus struct {
	ChainId         uint64 `json:"chain_id"`
	LastBlockHeight int64  `json:"last_block_height"`
	LastBlockHash   string `json:"last_block_hash"`
}

type ChainParams struct {
	ChainId   uint64 `json:"chain_id"`
	EpochSize *int64  `json:"epoch_size"`
}

type HeightMeta struct {
	Height int64  `json:"height"`
	Time   uint64 `json:"time"`

	Epoch       *int64 `json:"epoch"`
	LastInEpoch *bool  `json:"last_in_epoch"`
}

type Block struct {
	Height          int64      `json:"height"`
	Time            uint64     `json:"time"`
	Hash            string     `json:"hash"`
	ParentHash      string     `json:"parent_hash"`
	Coinbase        string     `json:"coinbase"`
	Root            string     `json:"root"`
	TxHash          string     `json:"tx_hash"`
	RecipientHash   string     `json:"recipient_hash"`
	Size            float64    `json:"size"`
	GasUsed         uint64     `json:"gas_used"`
	TotalDifficulty uint64     `json:"total_difficulty"`
	Extra           BlockExtra `json:"extra"`
	TxCount         int        `json:"tx_count"`
}

type BlockExtra struct {
	AddedValidators           []string                        `json:"added_validators"`
	AddedValidatorsPublicKeys []blscrypto.SerializedPublicKey `json:"added_validators_public_keys"`
	RemovedValidators         *big.Int                        `json:"removed_validators"`
	Seal                      []byte                          `json:"seal"`
	AggregatedSeal            IstanbulAggregatedSeal          `json:"aggregated_seal"`
	ParentAggregatedSeal      IstanbulAggregatedSeal          `json:"parent_aggregated_seal"`
}

type IstanbulAggregatedSeal struct {
	Bitmap    *big.Int `json:"bitmap"`
	Signature []byte   `json:"signature"`
	Round     uint64   `json:"round"`
}

type Transaction struct {
	Hash                string   `json:"hash"`
	To                  string   `json:"to"`
	Size                string   `json:"size"`
	Nonce               uint64   `json:"nonce"`
	GasPrice            *big.Int `json:"gas_price"`
	Gas                 uint64   `json:"gas"`
	GatewayFee          *big.Int `json:"gateway_fee"`
	GatewayFeeRecipient string   `json:"gateway_fee_recipient"`
	Index               uint     `json:"index"`
	GasUsed             uint64   `json:"gas_used"`
	CumulativeGasUsed   uint64   `json:"cumulative_gas_used"`
	Success             bool     `json:"success"`

	Operations []*Operation `json:"operations"`
}

type Operation struct {
	Name    string      `json:"name"`
	Details interface{} `json:"details"`
}

type Transfer struct {
	Index   uint64   `json:"index"`
	Type    string   `json:"type"`
	From    string   `json:"from"`
	To      string   `json:"to"`
	Value   *big.Int `json:"value"`
	Success bool     `json:"success"`
}

type ValidatorGroup struct {
	Index               uint64   `json:"index"`
	Address             string   `json:"address"`
	Name                string   `json:"name"`
	MetadataUrl         string   `json:"metadata_url"`
	Commission          *big.Int `json:"commission"`
	NextCommission      *big.Int `json:"next_commission"`
	NextCommissionBlock int64    `json:"next_commission_block"`
	SlashMultiplier     *big.Int `json:"slash_multiplier "`
	LastSlashed         *big.Int `json:"last_slashed"`
	ActiveVotes         *big.Int `json:"active_votes"`
	ActiveVotesUnits    *big.Int `json:"active_votes_units"`
	PendingVotes        *big.Int `json:"pending_votes"`
	Members             []string `json:"members"`
}

type Validator struct {
	Address        string   `json:"address"`
	Name           string   `json:"name"`
	MetadataUrl    string   `json:"metadata_url"`
	BlsPublicKey   []byte   `json:"bls_public_key"`
	EcdsaPublicKey []byte   `json:"ecdsa_public_key"`
	Signer         string   `json:"signer"`
	Affiliation    string   `json:"affiliation"`
	Score          *big.Int `json:"score"`
	Signed         *bool    `json:"signed"`
}

type AccountInfo struct {
	GoldBalance *big.Int `json:"gold_balance"`

	*Identity
	TotalLockedGold          *big.Int `json:"total_locked_gold"`
	TotalNonvotingLockedGold *big.Int `json:"total_nonvoting_locked_gold"`
	StableTokenBalance       *big.Int `json:"stable_token_balance"`
}

type Identity struct {
	Name        string `json:"name"`
	MetadataUrl string `json:"metadata_url"`
}

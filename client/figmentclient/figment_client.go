package figmentclient

import (
	"context"
	"fmt"
	kliento "github.com/celo-org/kliento/client"
	"github.com/celo-org/kliento/client/debug"
	"github.com/celo-org/kliento/registry"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	celoTypes "github.com/ethereum/go-ethereum/core/types"
	base "github.com/figment-networks/celo-indexer/client"
	"github.com/figment-networks/celo-indexer/utils"
	"math/big"
)

const (
	CeloClientFigment = "figment"
)

var (
	_ Client = (*client)(nil)
)

type Client interface {
	base.Client

	GetChainStatus(context.Context) (*ChainStatus, error)
	GetMetaByHeight(context.Context, int64) (*HeightMeta, error)
	GetBlockByHeight(context.Context, int64) (*Block, error)
	GetTransactionsByHeight(context.Context, int64) ([]*Transaction, error)
	GetValidatorGroupsByHeight(context.Context, int64) ([]*ValidatorGroup, error)
	GetValidatorsByHeight(context.Context, int64) ([]*Validator, error)
	GetAccountDetailsByAddressAndHeight(context.Context, string, int64) (*AccountDetails, error)
}

type client struct {
	cc                *kliento.CeloClient
	contractsRegistry *contractsRegistry
}

func New(url string) (*client, error) {
	cc, err := kliento.Dial(url)
	if err != nil {
		return nil, err
	}

	contractsRegistry, err := NewContractsRegistry(cc)
	if err != nil {
		return nil, err
	}

	return &client{
		cc:                cc,
		contractsRegistry: contractsRegistry,
	}, nil
}

func (l *client) GetName() string {
	return CeloClientFigment
}

func (l *client) Close() {
	l.cc.Close()
}

func (l *client) GetChainStatus(ctx context.Context) (*ChainStatus, error) {
	chainId, err := l.cc.Net.ChainId(ctx)
	if err != nil {
		return nil, err
	}

	last, err := l.cc.Eth.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	chain := &ChainStatus{
		ChainId:         chainId.Uint64(),
		LastBlockHeight: last.Number.Int64(),
		LastBlockHash:   last.Hash().String(),
	}

	return chain, nil
}

func (l *client) GetMetaByHeight(ctx context.Context, h int64) (*HeightMeta, error) {
	height := big.NewInt(h)
	l.contractsRegistry.setupContractsForHeight(ctx, height)

	heightMeta := &HeightMeta{
		Height: h,
	}

	rawBlock, err := l.cc.Eth.BlockByNumber(ctx, height)
	if err != nil {
		return nil, err
	}
	heightMeta.Time = rawBlock.Time()

	if l.contractsRegistry.validatorsContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		epoch, err := l.contractsRegistry.validatorsContract.GetEpochNumberOfBlock(opts, height)
		if err != nil {
			return nil, err
		}
		e := epoch.Int64()
		heightMeta.Epoch = &e
	}

	if l.contractsRegistry.electionContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		epochSize, err := l.contractsRegistry.electionContract.GetEpochSize(opts)
		if err != nil {
			return nil, err
		}
		isLastInEpoch := istanbul.IsLastBlockOfEpoch(height.Uint64(), epochSize.Uint64())
		heightMeta.LastInEpoch = &isLastInEpoch
	}

	return heightMeta, nil
}

func (l *client) GetBlockByHeight(ctx context.Context, h int64) (*Block, error) {
	height := big.NewInt(h)
	rawBlock, err := l.cc.Eth.BlockByNumber(ctx, height)
	if err != nil {
		return nil, err
	}

	nextHeight := big.NewInt(0)
	nextHeight.Add(height, big.NewInt(1))
	nextBlockHeader, err := l.cc.Eth.HeaderByNumber(ctx, nextHeight)
	if err != nil {
		return nil, err
	}

	extra, err := celoTypes.ExtractIstanbulExtra(nextBlockHeader)
	if err != nil {
		return nil, err
	}

	block := &Block{
		Hash:            rawBlock.Hash().String(),
		Height:          rawBlock.Number().Int64(),
		ParentHash:      rawBlock.ParentHash().String(),
		Time:            rawBlock.Time(),
		Size:            float64(rawBlock.Size()),
		GasUsed:         rawBlock.GasUsed(),
		Coinbase:        rawBlock.Coinbase().String(),
		Root:            rawBlock.Root().String(),
		TxHash:          rawBlock.TxHash().String(),
		RecipientHash:   rawBlock.ReceiptHash().String(),
		TotalDifficulty: rawBlock.TotalDifficulty().Uint64(),
		Extra: BlockExtra{
			AddedValidators:           utils.StringifyAddresses(extra.AddedValidators),
			AddedValidatorsPublicKeys: extra.AddedValidatorsPublicKeys,
			RemovedValidators:         extra.RemovedValidators,
			Seal:                      extra.Seal,
			AggregatedSeal: IstanbulAggregatedSeal{
				Bitmap:    extra.AggregatedSeal.Bitmap,
				Signature: extra.AggregatedSeal.Signature,
				Round:     extra.AggregatedSeal.Round.Uint64(),
			},
			ParentAggregatedSeal: IstanbulAggregatedSeal{
				Bitmap:    extra.ParentAggregatedSeal.Bitmap,
				Signature: extra.ParentAggregatedSeal.Signature,
				Round:     extra.ParentAggregatedSeal.Round.Uint64(),
			},
		},
		TxCount: len(rawBlock.Transactions()),
	}

	return block, nil
}

func (l *client) GetTransactionsByHeight(ctx context.Context, h int64) ([]*Transaction, error) {
	height := big.NewInt(h)
	block, err := l.cc.Eth.BlockByNumber(ctx, height)
	if err != nil {
		return nil, err
	}

	rawTransactions := block.Transactions()

	var transactions []*Transaction
	for _, tx := range rawTransactions {
		txHash := tx.Hash()

		receipt, err := l.cc.Eth.TransactionReceipt(ctx, txHash)
		if err != nil {
			return nil, err
		}

		operations, err := l.parseFromLogs(receipt.Logs)
		if err != nil {
			return nil, err
		}

		// Internal transfers
		internalTransfers, err := l.cc.Debug.TransactionTransfers(ctx, txHash)
		if err != nil {
			return nil, fmt.Errorf("can't run celo-rpc tx-tracer: %w", err)
		}
		operations = append(operations, l.parseFromInternalTransfers(internalTransfers)...)

		transaction := &Transaction{
			Nonce:      tx.Nonce(),
			GasPrice:   tx.GasPrice(),
			Gas:        tx.Gas(),
			GatewayFee: tx.GatewayFee(),

			Index:             receipt.TransactionIndex,
			GasUsed:           receipt.GasUsed,
			CumulativeGasUsed: receipt.CumulativeGasUsed,
			Success:           receipt.Status == celoTypes.ReceiptStatusSuccessful,
			Operations:        operations,
		}

		if tx.GatewayFeeRecipient() != nil {
			transaction.GatewayFeeRecipient = tx.GatewayFeeRecipient().String()
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (l *client) parseFromInternalTransfers(internalTransfers []debug.Transfer) []*Operation {
	var operations []*Operation
	for i, t := range internalTransfers {
		transfer := &Transfer{
			Index:   uint64(i),
			Type:    "transfer",
			From:    t.From.String(),
			To:      t.To.String(),
			Value:   t.Value,
			Success: t.Status == debug.TransferStatusSuccess,
		}

		operations = append(operations, &Operation{
			Name:    "internalTransfer",
			Details: transfer,
		})
	}
	return operations
}

func (l *client) parseFromLogs(logs []*celoTypes.Log) ([]*Operation, error) {
	var operations []*Operation
	for _, eventLog := range logs {
		if eventLog.Address == l.contractsRegistry.addresses[registry.ElectionContractID] {
			eventName, eventRaw, ok, err := l.contractsRegistry.electionContract.TryParseLog(*eventLog)
			if err != nil {
				return nil, fmt.Errorf("can't parse Election event: %w", err)
			}
			if !ok {
				continue
			}

			operations = append(operations, &Operation{
				Name:    eventName,
				Details: eventRaw,
			})
			// Election:
			//switch eventName {
			//case "ValidatorGroupVoteCast":
			//	// vote() [ValidatorGroupVoteCast] => lockNonVoting->lockVotingPending
			//	event := eventRaw.(*contracts.ElectionValidatorGroupVoteCast)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"group":   event.Group,
			//			"value":   event.Value,
			//		},
			//	}
			//	operations = append(operations, op)
			//case "ValidatorGroupVoteActivated":
			//	// activate() [ValidatorGroupVoteActivated] => lockVotingPending->lockVotingActive
			//	event := eventRaw.(*contracts.ElectionValidatorGroupVoteActivated)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"group":   event.Group,
			//			"value":   event.Value,
			//			"units":   event.Units,
			//		},
			//	}
			//	operations = append(operations, op)
			//case "ValidatorGroupPendingVoteRevoked":
			//	// revokePending() [ValidatorGroupPendingVoteRevoked] => lockVotingPending->lockNonVoting
			//	event := eventRaw.(*contracts.ElectionValidatorGroupPendingVoteRevoked)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"group":   event.Group,
			//			"value":   event.Value,
			//		},
			//	}
			//	operations = append(operations, op)
			//case "ValidatorGroupActiveVoteRevoked":
			//	// revokeActive() [ValidatorGroupActiveVoteRevoked] => lockVotingActive->lockNonVoting
			//	event := eventRaw.(*contracts.ElectionValidatorGroupActiveVoteRevoked)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"group":   event.Group,
			//			"value":   event.Value,
			//			"units":   event.Units,
			//		},
			//	}
			//	operations = append(operations, op)
			//}

		} else if eventLog.Address == l.contractsRegistry.addresses[registry.AccountsContractID] {
			eventName, eventRaw, ok, err := l.contractsRegistry.accountsContract.TryParseLog(*eventLog)
			if err != nil {
				return nil, fmt.Errorf("can't parse Accounts event: %w", err)
			}
			if !ok {
				continue
			}

			operations = append(operations, &Operation{
				Name:    eventName,
				Details: eventRaw,
			})
			// Accounts:
			//switch eventName {
			//case "AccountCreated":
			//	event := eventRaw.(*contracts.AccountsAccountCreated)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//		},
			//	}
			//	operations = append(operations, op)
			//case "VoteSignerAuthorized":
			//	event := eventRaw.(*contracts.AccountsVoteSignerAuthorized)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"signer":  event.Signer,
			//		},
			//	}
			//	operations = append(operations, op)
			//case "ValidatorSignerAuthorized":
			//	event := eventRaw.(*contracts.AccountsValidatorSignerAuthorized)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"signer":  event.Signer,
			//		},
			//	}
			//	operations = append(operations, op)
			//case "AttestationSignerAuthorized":
			//	event := eventRaw.(*contracts.AccountsAttestationSignerAuthorized)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"signer":  event.Signer,
			//		},
			//	}
			//	operations = append(operations, op)
			//}

		} else if eventLog.Address == l.contractsRegistry.addresses[registry.LockedGoldContractID] {
			eventName, eventRaw, ok, err := l.contractsRegistry.lockedGoldContract.TryParseLog(*eventLog)
			if err != nil {
				return nil, fmt.Errorf("can't parse LockedGold event: %w", err)
			}
			if !ok {
				continue
			}

			operations = append(operations, &Operation{
				Name:    eventName,
				Details: eventRaw,
			})

			//switch eventName {
			//case "GoldLocked":
			//	// lock() [GoldLocked + transfer] => main->lockNonVoting
			//	event := eventRaw.(*contracts.LockedGoldGoldLocked)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"value":  event.Value,
			//		},
			//	}
			//	operations = append(operations, op)
			//
			//case "GoldRelocked":
			//	// relock() [GoldRelocked] => lockPending->lockNonVoting
			//	event := eventRaw.(*contracts.LockedGoldGoldRelocked)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"value":  event.Value,
			//		},
			//	}
			//	operations = append(operations, op)
			//
			//case "GoldUnlocked":
			//	// unlock() [GoldUnlocked] => lockNonVoting->lockPending
			//	event := eventRaw.(*contracts.LockedGoldGoldUnlocked)
			//	op := &Operation{
			//		Details: map[string]interface{}{
			//			"account": event.Account,
			//			"value":  event.Value,
			//			"available": event.Available,
			//		},
			//	}
			//	operations = append(operations, op)
			//
			//case "GoldWithdrawn":
			//	// withdraw() [GoldWithdrawn + transfer] => lockPending->main
			//	event := eventRaw.(*contracts.LockedGoldGoldWithdrawn)
			//	operations = append(operations, *NewWithdrawGold(event.Account, lockedGoldAddr, event.Value, tobinTax))
			//
			//case "AccountSlashed":
			//	// slash() [AccountSlashed + transfer] => account:lockNonVoting -> beneficiary:lockNonVoting + governance:main
			//	event := eventRaw.(*contracts.LockedGoldAccountSlashed)
			//	operations = append(operations, *NewSlash(event.Slashed, event.Reporter, governanceAddr, lockedGoldAddr, event.Penalty, event.Reward, tobinTax))
			//
			//}
		} else if eventLog.Address == l.contractsRegistry.addresses[registry.StableTokenContractID] {
			eventName, eventRaw, ok, err := l.contractsRegistry.stableTokenContract.TryParseLog(*eventLog)
			if err != nil {
				return nil, fmt.Errorf("can't parse StableToken event: %w", err)
			}
			if !ok {
				continue
			}

			operations = append(operations, &Operation{
				Name:    eventName,
				Details: eventRaw,
			})
		}

	}
	return operations, nil
}

func (l *client) GetValidatorGroupsByHeight(ctx context.Context, h int64) ([]*ValidatorGroup, error) {
	height := big.NewInt(h)
	l.contractsRegistry.setupContractsForHeight(ctx, height)

	var validatorGroups []*ValidatorGroup

	if l.contractsRegistry.validatorsContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		rawValidatorGroups, err := l.contractsRegistry.validatorsContract.GetRegisteredValidatorGroups(opts)
		if err != nil {
			return nil, err
		}

		for i, rawValidatorGroup := range rawValidatorGroups {
			opts = &bind.CallOpts{Context: ctx}
			members, commission, nextCommission, nextCommissionBlock, _, slashMultiplier, lastSlashed, err := l.contractsRegistry.validatorsContract.GetValidatorGroup(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			opts = &bind.CallOpts{Context: ctx}
			activeVotes, err := l.contractsRegistry.electionContract.GetActiveVotesForGroup(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			opts = &bind.CallOpts{Context: ctx}
			activeVoteUnits, err := l.contractsRegistry.electionContract.GetActiveVoteUnitsForGroup(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			opts = &bind.CallOpts{Context: ctx}
			pendingVotes, err := l.contractsRegistry.electionContract.GetPendingVotesForGroup(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			validatorGroup := &ValidatorGroup{
				Index:               uint64(i),
				Address:             rawValidatorGroup.String(),
				Commission:          commission,
				NextCommission:      nextCommission,
				NextCommissionBlock: nextCommissionBlock.Int64(),
				SlashMultiplier:     slashMultiplier,
				LastSlashed:         lastSlashed,
				ActiveVotes:         activeVotes,
				ActiveVotesUnits:    activeVoteUnits,
				PendingVotes:        pendingVotes,
			}

			validatorGroup.Members = []string{}
			for _, member := range members {
				validatorGroup.Members = append(validatorGroup.Members, member.String())
			}

			validatorGroups = append(validatorGroups, validatorGroup)
		}
	}

	return validatorGroups, nil
}

func (l *client) GetValidatorsByHeight(ctx context.Context, h int64) ([]*Validator, error) {
	height := big.NewInt(h)
	l.contractsRegistry.setupContractsForHeight(ctx, height)

	var validators []*Validator

	if l.contractsRegistry.validatorsContract == nil {
		return validators, nil
	}

	opts := &bind.CallOpts{Context: ctx}
	rawValidators, err := l.contractsRegistry.validatorsContract.GetRegisteredValidators(opts)
	if err != nil {
		return nil, err
	}

	validationMap, err := l.getValidationMap(ctx, height)
	if err != nil {
		return nil, err
	}

	for _, rawValidator := range rawValidators {
		opts := &bind.CallOpts{Context: ctx}
		validatorDetails, err := l.contractsRegistry.validatorsContract.GetValidator(opts, rawValidator)
		if err != nil {
			return nil, err
		}

		validator := &Validator{
			Address:        rawValidator.String(),
			BlsPublicKey:   validatorDetails.BlsPublicKey,
			EcdsaPublicKey: validatorDetails.EcdsaPublicKey,
			Signer:         validatorDetails.Signer.String(),
			Affiliation:    validatorDetails.Affiliation.String(),
			Score:          validatorDetails.Score,
		}

		signed, ok := validationMap[rawValidator.String()]
		if ok {
			validator.Signed = &signed
		}

		validators = append(validators, validator)
	}

	return validators, nil
}

func (l *client) getValidationMap(ctx context.Context, height *big.Int) (map[string]bool, error) {
	validationMap := map[string]bool{}

	if l.contractsRegistry.electionContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		currentValidatorSigners, err := l.contractsRegistry.electionContract.GetCurrentValidatorSigners(opts)
		if err != nil {
			return nil, err
		}

		nextHeight := big.NewInt(0)
		nextHeight.Add(height, big.NewInt(1))
		nextBlockHeader, err := l.cc.Eth.HeaderByNumber(ctx, nextHeight)
		if err != nil {
			return nil, err
		}

		extra, err := celoTypes.ExtractIstanbulExtra(nextBlockHeader)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(currentValidatorSigners); i++ {
			signer := currentValidatorSigners[uint64(i)]
			signed := false

			if extra.ParentAggregatedSeal.Bitmap.Bit(i) == 1 {
				signed = true
			}

			validationMap[signer.String()] = signed
		}
	}

	return validationMap, nil
}

func (l *client) GetAccountDetailsByAddressAndHeight(ctx context.Context, rawAddress string, h int64) (*AccountDetails, error) {
	height := big.NewInt(h)
	l.contractsRegistry.setupContractsForHeight(ctx, height)

	address := common.HexToAddress(rawAddress)

	accountDetails := &AccountDetails{}

	goldAmount, err := l.cc.Eth.BalanceAt(ctx, address, height)
	if err != nil {
		return nil, err
	}
	accountDetails.GoldBalance = goldAmount

	if l.contractsRegistry.lockedGoldContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		totalLockedGold, err := l.contractsRegistry.lockedGoldContract.GetAccountTotalLockedGold(opts, address)
		if err != nil {
			return nil, err
		}
		accountDetails.TotalLockedGold = totalLockedGold

		opts = &bind.CallOpts{Context: ctx}
		totalNonvotingLockedGold, err := l.contractsRegistry.lockedGoldContract.GetAccountNonvotingLockedGold(opts, address)
		if err != nil {
			return nil, err
		}
		accountDetails.TotalNonvotingLockedGold = totalNonvotingLockedGold
	}

	if l.contractsRegistry.stableTokenContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		stableTokenBalance, err := l.contractsRegistry.stableTokenContract.BalanceOf(opts, address)
		if err != nil {
			return nil, err
		}
		accountDetails.StableTokenBalance = stableTokenBalance
	}

	return accountDetails, nil
}

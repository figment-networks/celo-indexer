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
	CeloClientFigment = "figment_celo_client"
)

var (
	_ Client = (*client)(nil)
)

type Client interface {
	base.Client

	GetChainStatus(context.Context) (*ChainStatus, error)
	GetChainParams(context.Context) (*ChainParams, error)
	GetMetaByHeight(context.Context, int64) (*HeightMeta, error)
	GetBlockByHeight(context.Context, int64) (*Block, error)
	GetTransactionsByHeight(context.Context, int64) ([]*Transaction, error)
	GetValidatorGroupsByHeight(context.Context, int64) ([]*ValidatorGroup, error)
	GetValidatorsByHeight(context.Context, int64) ([]*Validator, error)
	GetAccountByAddressAndHeight(context.Context, string, int64) (*AccountInfo, error)
	GetIdentityByHeight(context.Context, string, int64) (*Identity, error)
}

type client struct {
	cc *kliento.CeloClient
}

func New(url string) (*client, error) {
	cc, err := kliento.Dial(url)
	if err != nil {
		return nil, err
	}

	return &client{
		cc: cc,
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

func (l *client) GetChainParams(ctx context.Context) (*ChainParams, error) {
	chainParams := &ChainParams{}

	chainId, err := l.cc.Net.ChainId(ctx)
	if err != nil {
		return nil, err
	}
	chainParams.ChainId = chainId.Uint64()

	cr, err := NewContractsRegistry(l.cc, nil)
	if err != nil {
		return nil, err
	}
	setupErr := cr.setupContracts(ctx, registry.ElectionContractID)

	if cr.contractDeployed(registry.ElectionContractID) {
		opts := &bind.CallOpts{Context: ctx}
		epochSize, err := cr.electionContract.GetEpochSize(opts)
		if err != nil {
			return nil, err
		}
		e := epochSize.Int64()
		chainParams.EpochSize = &e
	}

	return chainParams, setupErr
}

func (l *client) GetMetaByHeight(ctx context.Context, h int64) (*HeightMeta, error) {
	height := big.NewInt(h)

	heightMeta := &HeightMeta{
		Height: h,
	}

	rawBlock, err := l.cc.Eth.BlockByNumber(ctx, height)
	if err != nil {
		return nil, err
	}
	heightMeta.Time = rawBlock.Time()

	cr, err := NewContractsRegistry(l.cc, height)
	if err != nil {
		return nil, err
	}
	setupErr := cr.setupContracts(ctx, registry.ValidatorsContractID, registry.ElectionContractID)

	if cr.contractDeployed(registry.ValidatorsContractID) {
		opts := &bind.CallOpts{Context: ctx}
		epoch, err := cr.validatorsContract.GetEpochNumberOfBlock(opts, height)
		if err != nil {
			return nil, err
		}
		e := epoch.Int64()
		heightMeta.Epoch = &e
	}

	if cr.contractDeployed(registry.ElectionContractID) {
		opts := &bind.CallOpts{Context: ctx}
		epochSize, err := cr.electionContract.GetEpochSize(opts)
		if err != nil {
			return nil, err
		}
		isLastInEpoch := istanbul.IsLastBlockOfEpoch(height.Uint64(), epochSize.Uint64())
		heightMeta.LastInEpoch = &isLastInEpoch
	}

	return heightMeta, setupErr
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

	cr, err := NewContractsRegistry(l.cc, height)
	if err != nil {
		return nil, err
	}
	setupErr := cr.setupContracts(ctx)

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

		var operations []*Operation

		// Internal transfers
		internalTransfers, err := l.cc.Debug.TransactionTransfers(ctx, txHash)
		if err != nil {
			return nil, fmt.Errorf("can't run celo-rpc tx-tracer: %w", err)
		}
		operations = append(operations, l.parseFromInternalTransfers(internalTransfers)...)

		// Operations from logs
		operationsFromLogs, err := l.parseFromLogs(cr, receipt.Logs)
		if err != nil {
			return nil, err
		}
		operations = append(operations, operationsFromLogs...)

		transaction := &Transaction{
			Hash:       tx.Hash().String(),
			Time:       block.Time(),
			Height:     block.Number().Int64(),
			Address:    tx.To().String(),
			Size:       tx.Size().String(),
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

		if tx.To() != nil {
			transaction.To = tx.To().String()
		}

		if tx.GatewayFeeRecipient() != nil {
			transaction.GatewayFeeRecipient = tx.GatewayFeeRecipient().String()
		}

		transactions = append(transactions, transaction)
	}

	return transactions, setupErr
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
			Name:    OperationTypeInternalTransfer,
			Details: transfer,
		})
	}
	return operations
}

func (l *client) parseFromLogs(cr *contractsRegistry, logs []*celoTypes.Log) ([]*Operation, error) {
	var operations []*Operation
	for _, eventLog := range logs {
		if eventLog.Address == cr.addresses[registry.ElectionContractID] && cr.contractDeployed(registry.ElectionContractID) {
			eventName, eventRaw, ok, err := cr.electionContract.TryParseLog(*eventLog)
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

		} else if eventLog.Address == cr.addresses[registry.AccountsContractID] && cr.contractDeployed(registry.AccountsContractID) {
			eventName, eventRaw, ok, err := cr.accountsContract.TryParseLog(*eventLog)
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

		} else if eventLog.Address == cr.addresses[registry.LockedGoldContractID] && cr.contractDeployed(registry.LockedGoldContractID) {
			eventName, eventRaw, ok, err := cr.lockedGoldContract.TryParseLog(*eventLog)
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

		} else if eventLog.Address == cr.addresses[registry.StableTokenContractID] && cr.contractDeployed(registry.StableTokenContractID) {
			eventName, eventRaw, ok, err := cr.stableTokenContract.TryParseLog(*eventLog)
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
		} else if eventLog.Address == cr.addresses[registry.GoldTokenContractID] && cr.contractDeployed(registry.GoldTokenContractID) {
			eventName, eventRaw, ok, err := cr.goldTokenContract.TryParseLog(*eventLog)
			if err != nil {
				return nil, fmt.Errorf("can't parse GoldToken event: %w", err)
			}
			if !ok {
				continue
			}

			operations = append(operations, &Operation{
				Name:    eventName,
				Details: eventRaw,
			})
		} else if eventLog.Address == cr.addresses[registry.ValidatorsContractID] && cr.contractDeployed(registry.ValidatorsContractID) {
			eventName, eventRaw, ok, err := cr.validatorsContract.TryParseLog(*eventLog)
			if err != nil {
				return nil, fmt.Errorf("can't parse Validators event: %w", err)
			}
			if !ok {
				continue
			}

			operations = append(operations, &Operation{
				Name:    eventName,
				Details: eventRaw,
			})
		} else if eventLog.Address == cr.addresses[registry.GovernanceContractID] && cr.contractDeployed(registry.GovernanceContractID) {
			eventName, eventRaw, ok, err := cr.governanceContract.TryParseLog(*eventLog)
			if err != nil {
				return nil, fmt.Errorf("can't parse Governance event: %w", err)
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

	cr, err := NewContractsRegistry(l.cc, height)
	if err != nil {
		return nil, err
	}
	err = cr.setupContracts(ctx, registry.ValidatorsContractID, registry.ElectionContractID, registry.AccountsContractID)
	if err != nil {
		return nil, err
	}

	var validatorGroups []*ValidatorGroup

	if cr.validatorsContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		rawValidatorGroups, err := cr.validatorsContract.GetRegisteredValidatorGroups(opts)
		if err != nil {
			return nil, err
		}

		for i, rawValidatorGroup := range rawValidatorGroups {
			opts = &bind.CallOpts{Context: ctx}
			members, commission, nextCommission, nextCommissionBlock, _, slashMultiplier, lastSlashed, err := cr.validatorsContract.GetValidatorGroup(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			opts = &bind.CallOpts{Context: ctx}
			activeVotes, err := cr.electionContract.GetActiveVotesForGroup(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			opts = &bind.CallOpts{Context: ctx}
			activeVoteUnits, err := cr.electionContract.GetActiveVoteUnitsForGroup(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			opts = &bind.CallOpts{Context: ctx}
			pendingVotes, err := cr.electionContract.GetPendingVotesForGroup(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			opts := &bind.CallOpts{Context: ctx}
			votingCap, err := cr.electionContract.GetNumVotesReceivable(opts, rawValidatorGroup)
			if err != nil {
				return nil, err
			}

			identity, err := l.getIdentity(ctx, cr, rawValidatorGroup.String())
			if err != nil {
				return nil, err
			}

			validatorGroup := &ValidatorGroup{
				Index:               uint64(i),
				Address:             rawValidatorGroup.String(),
				Name:                identity.Name,
				MetadataUrl:         identity.MetadataUrl,
				Commission:          commission,
				NextCommission:      nextCommission,
				NextCommissionBlock: nextCommissionBlock.Int64(),
				SlashMultiplier:     slashMultiplier,
				LastSlashed:         lastSlashed,
				ActiveVotes:         activeVotes,
				ActiveVotesUnits:    activeVoteUnits,
				PendingVotes:        pendingVotes,
				VotingCap:           votingCap,
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

	cr, err := NewContractsRegistry(l.cc, height)
	if err != nil {
		return nil, err
	}
	err = cr.setupContracts(ctx, registry.ValidatorsContractID, registry.ElectionContractID, registry.AccountsContractID)
	if err != nil {
		return nil, err
	}

	var validators []*Validator

	opts := &bind.CallOpts{Context: ctx}
	rawValidators, err := cr.validatorsContract.GetRegisteredValidators(opts)
	if err != nil {
		return nil, err
	}

	validationMap, err := l.getValidationMap(ctx, cr, height)
	if err != nil {
		return nil, err
	}

	for _, rawValidator := range rawValidators {
		opts := &bind.CallOpts{Context: ctx}
		validatorDetails, err := cr.validatorsContract.GetValidator(opts, rawValidator)
		if err != nil {
			return nil, err
		}

		identity, err := l.getIdentity(ctx, cr, rawValidator.String())
		if err != nil {
			return nil, err
		}

		validator := &Validator{
			Address:        rawValidator.String(),
			Name:           identity.Name,
			MetadataUrl:    identity.MetadataUrl,
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

func (l *client) getValidationMap(ctx context.Context, cr *contractsRegistry, height *big.Int) (map[string]bool, error) {
	validationMap := map[string]bool{}

	if cr.electionContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		currentValidatorSigners, err := cr.electionContract.GetCurrentValidatorSigners(opts)
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

func (l *client) GetAccountByAddressAndHeight(ctx context.Context, rawAddress string, h int64) (*AccountInfo, error) {
	height := big.NewInt(h)

	cr, err := NewContractsRegistry(l.cc, height)
	if err != nil {
		return nil, err
	}
	setupErr := cr.setupContracts(ctx, registry.AccountsContractID, registry.LockedGoldContractID, registry.StableTokenContractID)

	address := common.HexToAddress(rawAddress)

	accountInfo := &AccountInfo{}

	goldAmount, err := l.cc.Eth.BalanceAt(ctx, address, height)
	if err != nil {
		return nil, err
	}
	accountInfo.GoldBalance = goldAmount

	if cr.contractDeployed(registry.AccountsContractID) {
		identity, err := l.getIdentity(ctx, cr, rawAddress)
		if err != nil {
			return nil, err
		}
		accountInfo.Identity = identity
	}

	if cr.contractDeployed(registry.LockedGoldContractID) {
		opts := &bind.CallOpts{Context: ctx}
		totalLockedGold, err := cr.lockedGoldContract.GetAccountTotalLockedGold(opts, address)
		if err != nil {
			return nil, err
		}
		accountInfo.TotalLockedGold = totalLockedGold

		opts = &bind.CallOpts{Context: ctx}
		totalNonvotingLockedGold, err := cr.lockedGoldContract.GetAccountNonvotingLockedGold(opts, address)
		if err != nil {
			return nil, err
		}
		accountInfo.TotalNonvotingLockedGold = totalNonvotingLockedGold
	}

	if cr.contractDeployed(registry.StableTokenContractID) {
		opts := &bind.CallOpts{Context: ctx}
		stableTokenBalance, err := cr.stableTokenContract.BalanceOf(opts, address)
		if err != nil {
			return nil, err
		}
		accountInfo.StableTokenBalance = stableTokenBalance
	}

	return accountInfo, setupErr
}

func (l *client) GetIdentityByHeight(ctx context.Context, rawAddress string, h int64) (*Identity, error) {
	height := big.NewInt(h)

	cr, err := NewContractsRegistry(l.cc, height)
	if err != nil {
		return nil, err
	}
	err = cr.setupContracts(ctx, registry.AccountsContractID)
	if err != nil {
		return nil, err
	}

	return l.getIdentity(ctx, cr, rawAddress)
}

func (l *client) getIdentity(ctx context.Context, cr *contractsRegistry, rawAddress string) (*Identity, error) {
	address := common.HexToAddress(rawAddress)

	identity := &Identity{}
	if cr.accountsContract != nil {
		opts := &bind.CallOpts{Context: ctx}
		name, err := cr.accountsContract.GetName(opts, address)
		if err != nil {
			return nil, err
		}
		identity.Name = name

		metadataUrl, err := cr.accountsContract.GetMetadataURL(opts, address)
		if err != nil {
			return nil, err
		}
		identity.MetadataUrl = metadataUrl

	}

	return identity, nil
}

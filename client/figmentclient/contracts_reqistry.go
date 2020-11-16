package figmentclient

import (
	"context"
	kliento "github.com/celo-org/kliento/client"
	"github.com/celo-org/kliento/contracts"
	"github.com/celo-org/kliento/registry"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func NewContractsRegistry(cc *kliento.CeloClient) (*contractsRegistry, error) {
	reg, err := registry.New(cc)
	if err != nil {
		return nil, err
	}

	return &contractsRegistry{
		cc:        cc,
		reg:       reg,
		addresses: map[registry.ContractID]common.Address{},
	}, nil
}

type contractsRegistry struct {
	cc  *kliento.CeloClient
	reg registry.Registry

	addresses map[registry.ContractID]common.Address

	reserveContract      *contracts.Reserve
	stableTokenContract  *contracts.StableToken
	validatorsContract   *contracts.Validators
	lockedGoldContract   *contracts.LockedGold
	electionContract     *contracts.Election
	accountsContract     *contracts.Accounts
	goldTokenContract    *contracts.GoldToken
	chainParamsContract  *contracts.BlockchainParameters
	epochRewardsContract *contracts.EpochRewards
}

func (l *contractsRegistry) setupContractsForHeight(ctx context.Context, height *big.Int) {
	l.setupReserveContractForHeight(ctx, height)
	l.setupStableTokenContractForHeight(ctx, height)
	l.setupValidatorsContractForHeight(ctx, height)
	l.setupLockedGoldContractForHeight(ctx, height)
	l.setupElectionContractForHeight(ctx, height)
	l.setupAccountsContractForHeight(ctx, height)
	l.setupGoldTokenContractForHeight(ctx, height)
	l.setupEpochRewardsContractForHeight(ctx, height)
}

func (l *contractsRegistry) setupReserveContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.ReserveContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewReserve(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.ReserveContractID] = address
	l.reserveContract = contract

	return nil
}

func (l *contractsRegistry) setupStableTokenContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.StableTokenContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewStableToken(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.StableTokenContractID] = address
	l.stableTokenContract = contract

	return nil
}

func (l *contractsRegistry) setupValidatorsContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.ValidatorsContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewValidators(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.ValidatorsContractID] = address
	l.validatorsContract = contract

	return nil
}

func (l *contractsRegistry) setupLockedGoldContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.LockedGoldContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewLockedGold(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.LockedGoldContractID] = address
	l.lockedGoldContract = contract

	return nil
}

func (l *contractsRegistry) setupElectionContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.ElectionContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewElection(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.ElectionContractID] = address
	l.electionContract = contract

	return nil
}

func (l *contractsRegistry) setupAccountsContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.AccountsContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewAccounts(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.AccountsContractID] = address
	l.accountsContract = contract

	return nil
}

func (l *contractsRegistry) setupGoldTokenContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.GoldTokenContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewGoldToken(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.GoldTokenContractID] = address
	l.goldTokenContract = contract

	return nil
}

func (l *contractsRegistry) setupChainParamsContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.BlockchainParametersContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewBlockchainParameters(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.BlockchainParametersContractID] = address
	l.chainParamsContract = contract

	return nil
}

func (l *contractsRegistry) setupEpochRewardsContractForHeight(ctx context.Context, height *big.Int) error {
	address, err := l.reg.GetAddressFor(ctx, height, registry.EpochRewardsContractID)
	if err != nil {
		return err
	}
	contract, err := contracts.NewEpochRewards(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.BlockchainParametersContractID] = address
	l.epochRewardsContract = contract

	return nil
}

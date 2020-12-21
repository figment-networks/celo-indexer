package figmentclient

import (
	"context"
	"errors"
	kliento "github.com/celo-org/kliento/client"
	"github.com/celo-org/kliento/contracts"
	"github.com/celo-org/kliento/registry"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

var (
	ErrContractNotDeployed = errors.New("contract not deployed")
)

func NewContractsRegistry(cc *kliento.CeloClient, height *big.Int) (*contractsRegistry, error) {
	reg, err := registry.New(cc)
	if err != nil {
		return nil, err
	}

	return &contractsRegistry{
		cc:     cc,
		reg:    reg,
		height: height,

		addresses: map[registry.ContractID]common.Address{},
	}, nil
}

type contractsRegistry struct {
	cc     *kliento.CeloClient
	reg    registry.Registry
	height *big.Int

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
	governanceContract   *contracts.Governance
}

func contractIncluded(contracts []registry.ContractID, contractID registry.ContractID) bool {
	for _, id := range contracts {
		if id == contractID {
			return true
		}
	}
	return false
}

func (l *contractsRegistry) setupContracts(ctx context.Context, contracts ...registry.ContractID) error {
	if len(contracts) == 0 || contractIncluded(contracts, registry.ReserveContractID) {
		err := l.setupReserveContract(ctx)
		if err != nil {
			return err
		}
	}
	if len(contracts) == 0 || contractIncluded(contracts, registry.StableTokenContractID) {
		err := l.setupStableTokenContract(ctx)
		if err != nil {
			return err
		}
	}
	if len(contracts) == 0 || contractIncluded(contracts, registry.ValidatorsContractID) {
		err := l.setupValidatorsContract(ctx)
		if err != nil {
			return err
		}
	}
	if len(contracts) == 0 || contractIncluded(contracts, registry.LockedGoldContractID) {
		err := l.setupLockedGoldContract(ctx)
		if err != nil {
			return err
		}
	}
	if len(contracts) == 0 || contractIncluded(contracts, registry.ElectionContractID) {
		err := l.setupElectionContract(ctx)
		if err != nil {
			return err
		}
	}
	if len(contracts) == 0 || contractIncluded(contracts, registry.AccountsContractID) {
		err := l.setupAccountsContract(ctx)
		if err != nil {
			return err
		}
	}
	if len(contracts) == 0 || contractIncluded(contracts, registry.GoldTokenContractID) {
		err := l.setupGoldTokenContract(ctx)
		if err != nil {
			return err
		}
	}
	if len(contracts) == 0 || contractIncluded(contracts, registry.EpochRewardsContractID) {
		err := l.setupEpochRewardsContract(ctx)
		if err != nil {
			return err
		}
	}
	if len(contracts) == 0 || contractIncluded(contracts, registry.GovernanceContractID) {
		err := l.setupGovernanceContract(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) error {
	if err == kliento.ErrContractNotDeployed {
		return ErrContractNotDeployed
	} else {
		return err
	}
}

func (l *contractsRegistry) contractDeployed(contractId registry.ContractID) bool {
	_, ok := l.addresses[contractId]
	return ok
}

func (l *contractsRegistry) setupReserveContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.ReserveContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewReserve(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.ReserveContractID] = address
	l.reserveContract = contract

	return nil
}

func (l *contractsRegistry) setupStableTokenContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.StableTokenContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewStableToken(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.StableTokenContractID] = address
	l.stableTokenContract = contract

	return nil
}

func (l *contractsRegistry) setupValidatorsContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.ValidatorsContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewValidators(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.ValidatorsContractID] = address
	l.validatorsContract = contract

	return nil
}

func (l *contractsRegistry) setupLockedGoldContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.LockedGoldContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewLockedGold(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.LockedGoldContractID] = address
	l.lockedGoldContract = contract

	return nil
}

func (l *contractsRegistry) setupElectionContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.ElectionContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewElection(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.ElectionContractID] = address
	l.electionContract = contract

	return nil
}

func (l *contractsRegistry) setupAccountsContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.AccountsContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewAccounts(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.AccountsContractID] = address
	l.accountsContract = contract

	return nil
}

func (l *contractsRegistry) setupGoldTokenContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.GoldTokenContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewGoldToken(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.GoldTokenContractID] = address
	l.goldTokenContract = contract

	return nil
}

func (l *contractsRegistry) setupChainParamsContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.BlockchainParametersContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewBlockchainParameters(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.BlockchainParametersContractID] = address
	l.chainParamsContract = contract

	return nil
}

func (l *contractsRegistry) setupEpochRewardsContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.EpochRewardsContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewEpochRewards(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.BlockchainParametersContractID] = address
	l.epochRewardsContract = contract

	return nil
}

func (l *contractsRegistry) setupGovernanceContract(ctx context.Context) error {
	address, err := l.reg.GetAddressFor(ctx, l.height, registry.GovernanceContractID)
	if err != nil {
		return checkErr(err)
	}
	contract, err := contracts.NewGovernance(address, l.cc.Eth)
	if err != nil {
		return err
	}
	l.addresses[registry.GovernanceContractID] = address
	l.governanceContract = contract

	return nil
}

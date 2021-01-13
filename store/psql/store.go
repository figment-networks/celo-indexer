package psql

import (
	"github.com/figment-networks/celo-indexer/metric"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

const batchSize = 500

var (
	ErrNotFound = errors.New("record not found")
)

// NewIndexerMetric returns a new store from the connection string
func New(connStr string) (*Store, error) {
	conn, err := gorm.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	registerPlugins(conn)

	return &Store{
		db: conn,
	}, nil
}

// Store handles all database operations
type Store struct {
	db              *gorm.DB
	core            *core
	accounts        *accounts
	blocks          *blocks
	validators      *validators
	validatorGroups *validatorGroups
	governance      *governance
}

type core struct {
	*Database
	*Reports
	*Syncables
	*SystemEvents
}

type accounts struct {
	*AccountActivitySeq
}

type blocks struct {
	*BlockSeq
	*BlockSummary
}

type validators struct {
	*ValidatorAgg
	*ValidatorSeq
	*ValidatorSummary
}

type validatorGroups struct {
	*ValidatorGroupAgg
	*ValidatorGroupSeq
	*ValidatorGroupSummary
}

type governance struct {
	*ProposalAgg
	*GovernanceActivitySeq
}

// GetAccounts gets accounts
func (s *Store) GetAccounts() *accounts {
	if s.accounts == nil {
		s.accounts = &accounts{
			NewAccountActivitySeqStore(s.db),
		}
	}
	return s.accounts
}

// GetBlocks gets blocks
func (s *Store) GetBlocks() *blocks {
	if s.blocks == nil {
		s.blocks = &blocks{
			NewBlockSeqStore(s.db),
			NewBlockSummaryStore(s.db),
		}
	}
	return s.blocks
}

// GetDatabase gets database
func (s *Store) GetCore() *core {
	if s.core == nil {
		s.core = &core{
			NewDatabaseStore(s.db),
			NewReportsStore(s.db),
			NewSyncablesStore(s.db),
			NewSystemEventsStore(s.db),
		}
	}
	return s.core
}

// GetValidators gets validators
func (s *Store) GetValidators() *validators {
	if s.validators == nil {
		s.validators = &validators{
			NewValidatorAggStore(s.db),
			NewValidatorSeqStore(s.db),
			NewValidatorSummaryStore(s.db),
		}
	}
	return s.validators
}

// GetValidatorGroups gets validator groups
func (s *Store) GetValidatorGroups() *validatorGroups {
	if s.validatorGroups == nil {
		s.validatorGroups = &validatorGroups{
			NewValidatorGroupAggStore(s.db),
			NewValidatorGroupSeqStore(s.db),
			NewValidatorGroupSummaryStore(s.db),
		}
	}
	return s.validatorGroups
}

// GetGovernance gets governance
func (s *Store) GetGovernance() *governance {
	if s.governance == nil {
		s.governance = &governance{
			NewProposalAggStore(s.db),
			NewGovernanceActivitySeqStore(s.db),
		}
	}
	return s.governance
}

// Test checks the connection status
func (s *Store) Test() error {
	return s.db.DB().Ping()
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// SetDebugMode enabled detailed query logging
func (s *Store) SetDebugMode(enabled bool) {
	s.db.LogMode(enabled)
}

// registerPlugins registers gorm plugins
func registerPlugins(c *gorm.DB) {
	c.Callback().Create().Before("gorm:Create").Register("db_plugin:before_create", castQuantity)
	c.Callback().Update().Before("gorm:Update").Register("db_plugin:before_update", castQuantity)
}

// castQuantity casts decimal to quantity type
func castQuantity(scope *gorm.Scope) {
	for _, f := range scope.Fields() {
		v := f.Field.Type().String()
		if v == "types.Quantity" {
			f.IsNormal = true
			t := f.Field.Interface().(types.Quantity)
			f.Field = reflect.ValueOf(gorm.Expr("cast(? AS DECIMAL(65,0))", t.String()))
		}
	}
}

func LogQueryDuration(start time.Time, queryName string) {
	elapsed := time.Since(start)
	metric.DatabaseQueryDuration.WithLabelValues(queryName).Set(elapsed.Seconds())
}

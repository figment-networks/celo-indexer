package psql

const (
	bulkInsertGovernanceActivitySeqs = `
		INSERT INTO governance_activity_sequences (
		  height,
		  time,
		  proposal_id,
		  account,
		  transaction_hash,
		  kind,
		  data
		)
		VALUES @values;
	`
)

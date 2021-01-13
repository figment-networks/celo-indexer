package psql

const (
	bulkInsertGovernanceActivitySeqs = `
		INSERT INTO governance_activity_sequences (
		  height,
		  time,
		  proposal_id,
		  transaction_hash,
		  kind,
		  data
		)
		VALUES @values
		
		ON CONFLICT (height, proposal_id) DO UPDATE
		SET
		  transaction_hash = excluded.transaction_hash,
		  kind = excluded.kind,
		  data = excluded.data;
	`
)

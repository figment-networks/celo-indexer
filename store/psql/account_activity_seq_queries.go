package psql

var (
	bulkInsertAccountActivitySeqs = `
		INSERT INTO account_activity_sequences (
		  height,
		  time,
		  transaction_hash,
		  address,
		  amount,
		  kind,
          data
		)
		VALUES @values
		
		ON CONFLICT (height, transaction_hash) DO UPDATE
		SET
		  address = excluded.address,
		  amount = excluded.amount,
		  kind = excluded.kind,
		  data = excluded.data;
	`
)

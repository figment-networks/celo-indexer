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
		VALUES @values;
	`
)

package psql

const (
	bulkInsertValidatorSeqs = `
		INSERT INTO validator_sequences (
		  height,
		  time,
		  address,
		  affiliation,
		  signed,
          score
		)
		VALUES @values
		
		ON CONFLICT (height, address) DO UPDATE
		SET
		  affiliation = excluded.affiliation,
		  signed = excluded.signed,
		  score = excluded.score;
	`

	summarizeValidatorsQuerySelect = `
		address,
		DATE_TRUNC(?, time) AS time_bucket,
	
		AVG(signed::INT) AS signed_avg,
		MAX(signed::INT) AS signed_max,
		MIN(signed::INT) AS signed_min,
		AVG(score) AS score_avg,
		MAX(score) AS score_max,
		MIN(score) AS score_min
	`

	joinedAggregateSelect = `
		validator_sequences.height,
		validator_sequences.time,
		validator_sequences.address,
		validator_sequences.affiliation,
		validator_sequences.signed,
		validator_sequences.score,
		validator_aggregates.recent_name as name,
		validator_aggregates.recent_metadata_url as metadata_url
	`
)

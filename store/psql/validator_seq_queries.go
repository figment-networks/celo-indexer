package psql

const (
	bulkInsertValidatorSeqs = `
		INSERT INTO validator_sequences (
		  height,
		  time,
		  address,
		  name,
		  metadata_url,
		  affiliation,
		  signed,
          score
		)
		VALUES @values
		
		ON CONFLICT (height, address) DO UPDATE
		SET
		  name = excluded.name,
		  metadata_url = excluded.metadata_url,
		  affiliation = excluded.affiliation,
		  signed = excluded.signed,
		  score = excluded.score;
	`

	summarizeValidatorsForEraQuerySelect = `
	address,
	DATE_TRUNC(?, time) AS time_bucket,

	AVG(signed::INT) AS signed_avg,
   	MAX(signed::INT) AS signed_max,
   	MIN(signed::INT) AS signed_min,
   	AVG(score) AS score_avg,
   	MAX(score) AS score_max,
   	MIN(score) AS score_min
`
)

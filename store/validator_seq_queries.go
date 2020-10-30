package store

const (
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

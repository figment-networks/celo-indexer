package store

const (
	summarizeValidatorsForSessionQuerySelect = `
	address,
	DATE_TRUNC(?, time) AS time_bucket,
   	AVG(commission) AS commission_avg,
   	MAX(commission) AS commission_max,
   	MIN(commission) AS commission_min,
	AVG(active_votes) AS active_votes_avg,
   	MAX(active_votes) AS active_votes_max,
   	MIN(active_votes) AS active_votes_min,
	AVG(active_vote_units) AS active_vote_units_avg,
   	MAX(active_vote_units) AS active_vote_units_max,
   	MIN(active_vote_units) AS active_vote_units_min,
	AVG(pending_votes) AS pending_votes_avg,
   	MAX(pending_votes) AS pending_votes_max,
   	MIN(pending_votes) AS pending_votes_min
`
)

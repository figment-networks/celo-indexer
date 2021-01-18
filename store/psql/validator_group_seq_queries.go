package psql

const (
	bulkInsertValidatorGroupSeqs = `
		INSERT INTO validator_group_sequences (
		  height,
		  time,
		  address,
		  commission,
		  active_votes,
		  pending_votes,
		  voting_cap,
		  members_count,
          members_avg_signed
		)
		VALUES @values
		
		ON CONFLICT (height, address) DO UPDATE
		SET
		  commission = excluded.commission,
		  active_votes = excluded.active_votes,
		  pending_votes = excluded.pending_votes,
		  voting_cap = excluded.voting_cap,
		  members_count = excluded.members_count,
		  members_avg_signed = excluded.members_avg_signed;
	`

	summarizeValidatorGroupsQuerySelect = `
	address,
	DATE_TRUNC(?, time) AS time_bucket,
   	AVG(commission) AS commission_avg,
   	MAX(commission) AS commission_max,
   	MIN(commission) AS commission_min,
	AVG(active_votes) AS active_votes_avg,
   	MAX(active_votes) AS active_votes_max,
   	MIN(active_votes) AS active_votes_min,
	AVG(pending_votes) AS pending_votes_avg,
   	MAX(pending_votes) AS pending_votes_max,
   	MIN(pending_votes) AS pending_votes_min
`
)

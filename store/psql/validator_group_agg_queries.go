package psql

const (
	UpdateCalculate = `
		UPDATE validator_group_aggregates
		  	SET accumulated_uptime = t.members_count, 
				accumulated_uptime_count = t.accumulated_uptime_count
		FROM ( 
			SELECT 
				SUM(members_count) as members_count , 
				SUM(members_avg_signed) as accumulated_uptime_count
			FROM validator_group_sequences
				WHERE  address = ? 
		) t 
		WHERE 
		   address = ?
	`
)

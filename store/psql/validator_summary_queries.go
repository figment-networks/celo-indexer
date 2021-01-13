package psql

const (
	validatorSummaryForIntervalQuery = `
		SELECT * 
		FROM validator_summary 
		WHERE time_bucket >= (
			SELECT time_bucket 
			FROM validator_summary 
			WHERE time_interval = ?
			ORDER BY time_bucket DESC
			LIMIT 1
		) - ?::INTERVAL
			AND address = ? AND time_interval = ?
		ORDER BY time_bucket
	`

	allValidatorsSummaryForIntervalQuery = `
		SELECT
		  time_bucket,
		  time_interval,
		
		  AVG(signed_avg) AS signed_avg,
		  AVG(score_avg) AS score_avg,
		  MIN(score_min) AS score_min,
		  MAX(score_max) AS score_max
		FROM validator_summary
		WHERE time_bucket >= (
			SELECT time_bucket 
			FROM validator_summary 
			WHERE time_interval = ?
			ORDER BY time_bucket DESC 
			LIMIT 1
		) - ?::INTERVAL
			AND time_interval = ?
		GROUP BY time_bucket, time_interval
		ORDER BY time_bucket
	`

	validatorSummaryActivityPeriodsQuery = `
		WITH cte AS (
			SELECT
			  time_bucket,
			  sum(CASE WHEN diff IS NULL OR diff > ? :: INTERVAL
				THEN 1
				  ELSE NULL END)
			  OVER (
				ORDER BY time_bucket ) AS period
			FROM (
				   SELECT
					 time_bucket,
					 time_bucket - lag(time_bucket, 1)
					 OVER (
					   ORDER BY time_bucket ) AS diff
				   FROM validator_summary
				   WHERE time_interval = ? AND index_version = ?
				 ) AS x
		)
		SELECT
		  period,
		  MIN(time_bucket),
		  MAX(time_bucket)
		FROM cte
		GROUP BY period
		ORDER BY period
	`
)

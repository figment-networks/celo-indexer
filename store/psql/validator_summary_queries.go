package psql

const (
	bulkInsertValidatorSummaries = `
		INSERT INTO validator_summary (
	      	time_interval,
			time_bucket,
			index_version,
			address,
			score_avg,
			score_max,
			score_min,
			signed_avg,
			signed_min,
			signed_max
		)
		VALUES @values
		
		ON CONFLICT (time_interval, time_bucket, index_version, address) DO UPDATE
		SET
		  score_avg = excluded.score_avg,
		  score_max = excluded.score_max,
		  score_min = excluded.score_min,
		  signed_avg = excluded.signed_avg,
		  signed_min = excluded.signed_min,
		  signed_max = excluded.signed_max;
	`

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

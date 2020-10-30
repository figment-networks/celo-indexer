package store

const (
	validatorGroupSummaryForIntervalQuery = `
SELECT * 
FROM validator_group_summary 
WHERE time_bucket >= (
	SELECT time_bucket 
	FROM validator_group_summary 
	WHERE time_interval = ?
	ORDER BY time_bucket DESC
	LIMIT 1
) - ?::INTERVAL
	AND address = ? AND time_interval = ?
ORDER BY time_bucket
`

	allValidatorGroupsSummaryForIntervalQuery = `
SELECT
  time_bucket,
  time_interval,

  AVG(commission_avg) AS commission_avg,
  MIN(commission_min) AS commission_min,
  MAX(commission_max) AS commission_max,
  AVG(active_votes_avg) AS active_votes_avg,
  MIN(active_votes_min) AS active_votes_min,
  MAX(active_votes_max) AS active_votes_max,
  AVG(active_vote_units_avg) AS active_vote_units_avg,
  MIN(active_vote_units_min) AS active_vote_units_min,
  MAX(active_vote_units_max) AS active_votes_max,
  AVG(pending_votes_avg) AS pending_votes_avg,
  MIN(pending_votes_min) AS pending_votes_min,
  MAX(pending_votes_max) AS pending_votes_max
FROM validator_group_summary
WHERE time_bucket >= (
	SELECT time_bucket 
	FROM validator_group_summary 
	WHERE time_interval = ?
	ORDER BY time_bucket DESC 
	LIMIT 1
) - ?::INTERVAL
	AND time_interval = ?
GROUP BY time_bucket, time_interval
ORDER BY time_bucket
`

	validatorGroupSummaryActivityPeriodsQuery = `
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
           FROM validator_group_summary
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

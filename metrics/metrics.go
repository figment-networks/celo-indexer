package metrics

import (
	"github.com/figment-networks/indexing-engine/metrics"
)

var (
	PipelineUsecaseDuration = metrics.MustNewGaugeWithTags(metrics.Options{
		Namespace: "indexer",
		Subsystem: "pipeline",
		Name:      "usecase_duration",
		Desc:      "The total time spent executing a usecase",
		Tags:      []string{"task"},
	})

	PipelineDatabaseSizeAfterHeight = metrics.MustNewGaugeWithTags(metrics.Options{
		Namespace: "indexer",
		Subsystem: "pipeline",
		Name:      "database_size",
		Desc:      "The size of the database after indexing a height",
	})

	PipelineRequestCountAfterHeight = metrics.MustNewGaugeWithTags(metrics.Options{
		Namespace: "indexer",
		Subsystem: "pipeline",
		Name:      "request_count",
		Desc:      "The total number of requests made for one height",
	})

	DatabaseQueryDuration = metrics.MustNewGaugeWithTags(metrics.Options{
		Namespace: "indexer",
		Subsystem: "database",
		Name:      "query_duration",
		Desc:      "The total time required to execute query on database",
		Tags:      []string{"query"},
	})

	ServerRequestDuration = metrics.MustNewHistogramWithTags(metrics.HistogramOptions{
		Namespace: "indexer",
		Subsystem: "server",
		Name:      "request_duration",
		Desc:      "The total time spent handling an HTTP request",
		Tags:      []string{"request"},
	})
)

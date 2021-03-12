package metrics

import "time"

func LogUsecaseDuration(start time.Time, useCaseName string) {
	elapsed := time.Since(start)
	PipelineUsecaseDuration.WithLabels(useCaseName).Set(elapsed.Seconds())
}

func LogQueryDuration(start time.Time, queryName string) {
	elapsed := time.Since(start)
	DatabaseQueryDuration.WithLabels(queryName).Set(elapsed.Seconds())
}

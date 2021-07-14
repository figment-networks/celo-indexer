package psql

const (
	bulkInsertJobs = "INSERT INTO jobs (height, created_at, updated_at) VALUES @values"
)

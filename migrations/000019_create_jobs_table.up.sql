CREATE TABLE jobs (
  id          BIGSERIAL NOT NULL PRIMARY KEY,
  height      INTEGER NOT NULL,
  run_count   INTEGER NOT NULL DEFAULT 0,
  last_error  TEXT,
  started_at  TIMESTAMP WITH TIME ZONE,
  finished_at TIMESTAMP WITH TIME ZONE,

  created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at  TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX index_jobs_on_height ON jobs (height);

CREATE INDEX index_jobs_on_finished_at ON jobs (finished_at);

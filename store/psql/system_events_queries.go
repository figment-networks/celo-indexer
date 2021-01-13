package psql

var (
	bulkInsertSystemEvents = `
		INSERT INTO system_events (
		  height,
		  time,
		  actor,
		  kind,
          data
		)
		VALUES @values
		
		ON CONFLICT (height, actor, kind) DO UPDATE
		SET
		  actor = excluded.actor,
		  kind = excluded.kind,
		  data = excluded.data;
	`
)

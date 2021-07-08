### Description

Celo Indexer is responsible for fetching and indexing Celo data.

### External Packages:
* `kliento` - Celo client
* `indexing-engine` - A backbone for indexing process
* `gin` - Http server
* `gorm` - ORM with PostgreSQL interface
* `cron` - Cron jobs runner
* `zap` - logging

### Environmental variables:

* `APP_ENV` - application environment (development | production)
* `NODE_URL` - url to celo node
* `SERVER_ADDR` - address to use for API
* `SERVER_PORT` - port to use for API
* `FIRST_BLOCK_HEIGHT` - height of first block in chain
* `INDEX_WORKER_INTERVAL` - index interval for worker
* `SUMMARIZE_WORKER_INTERVAL` - summary interval for worker
* `PURGE_WORKER_INTERVAL` - purge interval for worker
* `DEFAULT_BATCH_SIZE` - syncing batch size. Setting this value to 0 means no batch size
* `DATABASE_DSN` - PostgreSQL database URL
* `DEBUG` - turn on db debugging mode
* `LOG_LEVEL` - level of log
* `LOG_OUTPUT` - log output (ie. stdout or /tmp/logs.json)
* `ROLLBAR_ACCESS_TOKEN` - Rollbar access token for error reporting
* `ROLLBAR_SERVER_ROOT` - Rollbar server root for error reporting
* `INDEXER_METRIC_ADDR` - Prometheus server address for indexer metrics
* `SERVER_METRIC_ADDR` - Prometheus server address for server metrics
* `METRIC_SERVER_URL` - Url at which metrics will be accessible (for both indexer and server)
* `PURGE_BLOCK_INTERVAL` - Block sequence older than given interval will be purged
* `PURGE_BLOCK_HOURLY_SUMMARY_INTERVAL` - Block hourly summary records older than given interval will be purged
* `PURGE_BLOCK_DAILY_SUMMARY_INTERVAL` - Block daily summary records older than given interval will be purged
* `PURGE_VALIDATOR_INTERVAL` - Validator sequence older than given interval will be purged
* `PURGE_VALIDATOR_HOURLY_SUMMARY_INTERVAL` - Validator hourly summary records older than given interval will be purged
* `PURGE_VALIDATOR_DAILY_SUMMARY_INTERVAL` - Validator daily summary records older than given interval will be purged
* `INDEXER_TARGETS_FILE` - JSON file with targets and its task names
* `FETCH_WORKERS` - Space-separated list of fetch worker endpoints
* `FETCH_WORKER_ADDR` - Fetch worker address
* `FETCH_WORKER_PORT` - Fetch worker port
* `FETCH_INTERVAL` - Processing interval for the fetch manager
* `REDIS_URL` - Redis server URL
* `REDIS_EXP` - Expiration time for the data stored in Redis

### Available endpoints:

| Method | Path                                 | Description                                                 | Params                                                                                                                                                |
|--------|------------------------------------  |-------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| GET    | `/health`                            | health endpoint                                             | -                                                                                                                                                     |
| GET    | `/status`                            | status of the application and chain                         | include_chain (bool, optional) -   when true, returns chain status                                                                                                                                             |
| GET    | `/block`                             | return block by height                                      | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/block_times/:limit`                | get last x block times                                      | limit (required) - limit of blocks                                                                                                                    |
| GET    | `/blocks_summary`                    | get block summary                                           | interval (required) - time interval [hourly or daily] period (required) - summary period [ie. 24 hours]                                               |
| GET    | `/transactions`                      | get list of transactions                                    | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/account/:address`                  | get account information for height                          | address (required) - address  height (optional) - height [Default: 0 = last]                                                                  |
| GET    | `/account_details/:address`          | get account details                                         | address (required) - address      limit (required) - number of recent account activities                                                                                                            |
| GET    | `/validators`                        | get list of validators                                      | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/validators/for_min_height/:height` | get the list of validators for height greater than provided | height (required) - height [Default: 0 = last]                                                                                                        |
| GET    | `/validator/:address`                | get validator by address                                    | address (required) - validator's address    sequences_limit (required) - number of last sequences to include    eras_limit (required) - number of last eras to include                                                                                                      |
| GET    | `/validators_summary`                | validator summary                                           | interval (required) - time interval [hourly or daily] period (required) - summary period [ie. 24 hours]  address (optional) - validator's address |
| GET    | `/validator_groups`                  | get list of validator groups                                | height (optional) - height [Default: 0 = last]                                                                                                        |
| GET    | `/validator_group/:address`          | get validator group by address                              | address (required) - validator's address    sequences_limit (required) - number of last sequences to include    eras_limit (required) - number of last eras to include                                                                                                      |
| GET    | `/validator_groups_summary`          | validator group summary                                     | interval (required) - time interval [hourly or daily] period (required) - summary period [ie. 24 hours]  address (optional) - validator's address |
| GET    | `/system_events/:address`            | system events for given actor                               | `address (required)` - address of account `after (optional)` - return events after with height greater than provided height  `kind (optional)` - system event kind |
| GET    | `/proposals`                         | get list of all proposals                                   | `cursor (optional)` - paging cursor `page_size (optional)` - size of one page of results |
| GET    | `/proposals/:proposal_id/activity`   | get governance activity on given proposal                   | `proposal_id (required)` - ID of proposal `cursor (optional)` - paging cursor  `page_size (optional)` - size of one page of results |

### Running app

Once you have created a database and specified all configuration options, you
need to migrate the database. You can do that by running the command below:

```bash
celo-indexer -config path/to/config.json -cmd=migrate
```

Start the data fetcher:

```bash
celo-indexer -config path/to/config.json -cmd=fetch_worker
celo-indexer -config path/to/config.json -cmd=fetch_manager
```

Start the indexer:

```bash
celo-indexer -config path/to/config.json -cmd=worker
```

Start the API server:

```bash
celo-indexer -config path/to/config.json -cmd=server
```

### Running one-off commands

Start indexer:
```bash
celo-indexer -config path/to/config.json -cmd=indexer_start
```

Create summary tables for sequences:
```bash
celo-indexer -config path/to/config.json -cmd=indexer_summarize
```

Purge old data:
```bash
celo-indexer -config path/to/config.json -cmd=indexer_purge
```

### Running tests

To run tests with coverage you can use `test` Makefile target:
```shell script
make test
```

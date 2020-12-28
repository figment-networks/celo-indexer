module github.com/figment-networks/celo-indexer

go 1.14

// Use this to use external build
replace github.com/ethereum/go-ethereum => github.com/celo-org/celo-blockchain v1.1.0

require (
	github.com/celo-org/kliento v0.1.2-0.20200608140637-c5afc8cf0f44
	github.com/ethereum/go-ethereum v1.9.8
	github.com/figment-networks/indexing-engine v0.1.11
	github.com/gin-gonic/gin v1.5.0
	github.com/golang-migrate/migrate/v4 v4.11.0
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.4.2
	github.com/jinzhu/gorm v1.9.12
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.7.1
	github.com/prometheus/common v0.10.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/rollbar/rollbar-go v1.2.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20200305140159-d7d444866696 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
)

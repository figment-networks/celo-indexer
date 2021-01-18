package config

import "fmt"

const (
	AppName    = "celo-indexer"
	AppVersion = "0.1.6"
	GitCommit  = "-"
	GoVersion  = "1.14"
)

func VersionString() string {
	return fmt.Sprintf(
		"%s %s (git: %s, %s)",
		AppName,
		AppVersion,
		GitCommit,
		GoVersion,
	)
}

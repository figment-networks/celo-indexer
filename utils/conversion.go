package utils

import (
	"github.com/ethereum/go-ethereum/common"
)

func StringifyAddresses(addresses []common.Address) []string {
	var results []string
	for _, address := range addresses {
		results = append(results, address.String())
	}
	return results
}
package staking

import (
	"log"
	"math/big"
	"os"

	bootokenContract "github.com/dickmanben/sol_go/bootoken"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func fetchTotalStakedByAddress(contractAddress string) (map[string]int, error) {
	// Define the variable holding this information
	addressTotals := map[string]int{}

	// Setup an ethereum client, I'm using infura
	client, err := ethclient.Dial(os.Getenv("ETH_CLIENT_ADDRESS"))
	if err != nil {
		return nil, err
	}

	// Load the contract into the ERC721Enum struct generated from solidity
	instance, err := bootokenContract.NewBootoken(common.HexToAddress(contractAddress), client)
	if err != nil {
		return nil, err
	}

	// Loop through total supply
	for i := 0; i < 12345; i++ {
		if i%1000 == 0 {
			log.Printf("Received %v\n", i)
		}
		// Call the ownerOf method with the item number to get the
		// address that owns that NFT
		add, err := instance.TokenOwners(nil, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}
		if _, ok := addressTotals[add.String()]; !ok {
			// add to map
			addressTotals[add.String()] = 1
		} else {
			// add to map
			addressTotals[add.String()] = addressTotals[add.String()] + 1
		}
	}
	return addressTotals, nil
}

//
// ScrapeStakers gets the owner totals and outputs
// them to a json file
//
func ScrapeStakers(address string) (map[int]int, error) {
	aggregatedOwners, err := fetchTotalStakedByAddress(address)
	if err != nil {
		return nil, err
	}
	groupedTotals := map[int]int{}

	for _, total := range aggregatedOwners {
		if total == 0 {
			continue
		}
		if _, ok := groupedTotals[total]; !ok {
			groupedTotals[total] = 1
		} else {
			groupedTotals[total] = groupedTotals[total] + 1
		}
	}
	return groupedTotals, nil
}

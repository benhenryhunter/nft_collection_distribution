package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	contract "github.com/dickmanben/sol_go/erc721enum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//
// AggregateContractOwners iterates through an ERC721Enum address
// and returns a map of addresses with their total number of items
// from the collection
//
func AggregateContractOwners(address string) map[string]interface{} {
	// Define the variable holding this information
	addressTotals := map[string]interface{}{}

	// Setup an ethereum client, I'm using infura
	client, err := ethclient.Dial(os.Getenv("ETH_CLIENT_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}

	// Convert the string address to proper type
	contractAddress := common.HexToAddress(address)

	// Load the contract into the ERC721Enum struct generated from solidity
	instance, err := contract.NewMain(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// Call the totalSupply method on the contract to get the amount
	// of NFTs available
	total, err := instance.TotalSupply(nil)
	if err != nil {
		panic(err)
	}

	// Convert the bigInt into an int that makes comparisons easier
	totalSupply := int(total.Int64())

	// Loop through total supply
	for i := 0; i < totalSupply; i++ {
		if i%1000 == 0 {
			log.Printf("Received %v\n", i)
		}
		// Call the ownerOf method with the item number to get the
		// address that owns that NFT
		add, err := instance.OwnerOf(nil, big.NewInt(int64(i)))
		if err != nil {
			panic(err)
		}
		if _, ok := addressTotals[add.String()]; !ok {
			// add to map
			addressTotals[add.String()] = 1
		} else {
			// add to map
			addressTotals[add.String()] = addressTotals[add.String()].(int) + 1
		}
	}
	return addressTotals
}

//
// ownersToMap reads the scraped JSON file and groups
// addresses by how many they hold
//
func ownersToMap(fileName string) map[int]int {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	ownerTotals := map[string]int{}
	if err := json.Unmarshal(b, &ownerTotals); err != nil {
		panic(err)
	}

	groupedTotals := map[int]int{}

	for _, total := range ownerTotals {
		if _, ok := groupedTotals[total]; !ok {
			groupedTotals[total] = 1
		} else {
			groupedTotals[total] = groupedTotals[total] + 1
		}
	}

	groupedBytes, err := json.MarshalIndent(groupedTotals, "", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("groupedOwners.json", groupedBytes, 0644)
	if err != nil {
		panic(err)
	}
	return groupedTotals
}

//
// scrapeOwners gets the owner totals and outputs
// them to a json file
//
func scrapeOwners(address string) {
	aggregatedOwners := AggregateContractOwners(address)
	b, err := json.MarshalIndent(aggregatedOwners, "", " ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("ownerMap.json", b, 0644)
	if err != nil {
		panic(err)
	}
}

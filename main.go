package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/subosito/gotenv"
)

func init() {
	if os.Getenv("env") == "prod" {
		gotenv.Load("prod.env")
	} else {
		gotenv.Load()
	}
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

func main() {
	etherealsWTFContractAddress := "0xfc778be06c9a58f8f3e5e99216efbb28f750bc98"
	scrapeOwners(etherealsWTFContractAddress)
	// ownersToMap("ownerMap.json")
}

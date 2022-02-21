package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadMapFromFile(fileName string) map[string]int {
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
	return ownerTotals
}

package main

import (
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

func main() {
	// ScrapeStakers("0x69a96059cc35da280af8005d165da1d040297696")
	// groupedStakers := ownersToMap("groupedStakers.json")

}

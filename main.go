package main

import (
	"os"

	commands "github.com/yutaoj/cloud-mta-yutaoj/cmd"
)

func main() {
	// Execute CLI Root commands
	err := commands.Execute()
	if err != nil {
		os.Exit(1)
	}
}

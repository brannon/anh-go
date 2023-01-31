package main

import (
	"log"
	"os"

	"github.com/brannon/anh-go/cli/cmd"
)

func main() {
	err := cmd.Execute(os.Args)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}

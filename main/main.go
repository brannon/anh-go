package main

import (
	"log"
	"os"

	"github.com/brannon/anh-go/cmd"
)

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}

func run(args []string) error {
	return cmd.Execute(args)
}

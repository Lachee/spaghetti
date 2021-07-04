package main

import (
	"fmt"
	"log"
	"os"

	"github.com/lachee/noodle"
	"github.com/lachee/xve-graph/src/spaghetti"
)

func main() {
	args := os.Args
	fmt.Println(args)

	// Before we run, preprocess
	var spag = &spaghetti.Application{}
	initializeBuildTag(spag)

	// Run
	var app noodle.Application = spag
	exitCode := noodle.Run(app, args[0])
	log.Println("Exited with code", exitCode)
}

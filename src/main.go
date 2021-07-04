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

	var app noodle.Application
	app = &spaghetti.Application{}
	exitCode := noodle.Run(app, args[0])
	log.Println("Exited with code", exitCode)
}

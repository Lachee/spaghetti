// +build development

package main

import (
	"github.com/lachee/xve-graph/src/spaghetti"
	"log"
)

func initializeBuildTag(app *spaghetti.Application) {
	log.Println("Spaghetti Development Build")
	app.EnableDebugger()
}

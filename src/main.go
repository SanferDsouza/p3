// ~/~ begin <<docs/index.md#main>>[init]
package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	log.Println("starting p3")

	// ~/~ begin <<docs/index.md#setup-flags>>[init]
	configPath := flag.String("config", "", "path to config file")
	// ~/~ end
	flag.Parse()
	// ~/~ begin <<docs/index.md#validate-flags>>[init]
	if len(*configPath) == 0 {
		log.Fatalln("--config flag was provided but no value was given")
	}
	// ~/~ begin <<docs/index.md#validate-config-exists>>[init]
	if _, err := os.Stat(*configPath); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("config file '%s' does not exist: %v", *configPath, err)
		} else {
			log.Fatalf("could not stat '%s': %v", *configPath, err)
		}
	}
	// ~/~ end
	log.Printf("config file '%s' was provided\n", *configPath)
	// ~/~ end
// ~/~ end

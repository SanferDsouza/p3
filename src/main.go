// ~/~ begin <<docs/index.md#main>>[init]
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// ~/~ begin <<docs/index.md#hocon-imports>>[init]
	"github.com/gurkankaymak/hocon"
	// ~/~ end
)

// ~/~ begin <<docs/index.md#hash-kinds>>[init]
type HashKind int

const (
	nokind HashKind = -1
	sha256 HashKind = iota
)
// ~/~ end

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

	// ~/~ begin <<docs/index.md#hocon-parse>>[init]
	config, err := hocon.ParseResource(*configPath)
	if err != nil {
		log.Fatalf("could not parse hocon file: %v", err)
	}
	phrases := config.GetArray("phrases")
	if phrases == nil {
		log.Fatalln("malformed config file, cannot find 'phrases'")
	}
	for _, phrase := range phrases {
		configPhrase, err := hocon.ParseString(phrase.String())
		if err != nil {
			log.Fatalf("could not parse phrase element: %v", err)
		}

	    // ~/~ begin <<docs/index.md#parse-config-phrase>>[init]
	    hintDirty := configPhrase.Get("hint").String()
	    hint := strings.Trim(hintDirty, "\"")

	    kindWithHashDirty := configPhrase.Get("hash").String()
	    kindWithHash := strings.Trim(kindWithHashDirty, "\"")
	    // ~/~ end

		kind, hash, err := extractHash(kindWithHash)
		if err != nil {
			log.Fatalf("could not extract hash in '%s': %v", kindWithHash, err)
		}
	    log.Printf("found kind=%s, hash=%s, hint=%s", kind, hash, hint)
	}
	// ~/~ end
}

// ~/~ begin <<docs/index.md#helpers>>[init]
func extractHash(kindWithHash string) (HashKind, string, error) {
	before, after, found := strings.Cut(kindWithHash, "-")
	if !found {
		err := fmt.Errorf("could not find hashkind separator '-' in %s", kindWithHash)
		return nokind, before, err
	}
	var kind HashKind
	switch before {
	case "sha256":
		kind = sha256
	default:
		err := fmt.Errorf("could not determine hashkind of '%s'", before)
		return nokind, before, err
	}
	return kind, after, nil
}
// ~/~ end
// ~/~ end

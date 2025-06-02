// ~/~ begin <<docs/index.md#main>>[init]
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"

	// ~/~ begin <<docs/index.md#hocon-imports>>[init]
	"github.com/gurkankaymak/hocon"
	"golang.org/x/term"
	// ~/~ end
)

// ~/~ begin <<docs/index.md#hash-kinds>>[init]
type HashKind int

// ~/~ begin <<docs/index.md#hash-kinds-string>>[init]
func (hk HashKind) String() string {
	switch hk {
	case sha256Hk:
		return "sha256"
	case nokind:
		return "nokind (error)"
	default:
		s := fmt.Sprintf("cannot convert to string, unexpected HashKind %d", hk)
		panic(s)
	}
}

// ~/~ end

const (
	nokind   HashKind = -1
	sha256Hk HashKind = iota
)

// ~/~ end

// ~/~ begin <<docs/index.md#define-p3-phrases-type>>[init]
type P3Phrase struct {
	Hash string
	Hint string
	Kind HashKind
}

// ~/~ begin <<docs/index.md#p3-phrase-string-method>>[init]
func (p3p *P3Phrase) String() string {
	var sb strings.Builder
	sb.WriteString("hash=")
	sb.WriteString(p3p.Hash)
	sb.WriteString(",")
	sb.WriteString("hint=")
	sb.WriteString(p3p.Hint)
	sb.WriteString(",")
	sb.WriteString("kind=")
	sb.WriteString(p3p.Kind.String())
	return sb.String()
}

// ~/~ end
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

	// ~/~ begin <<docs/index.md#initialize-p3-phrases>>[init]
	p3Phrases := make([]P3Phrase, 0, len(phrases))
	// ~/~ end

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

		// ~/~ begin <<docs/index.md#add-to-p3-phrases>>[init]
		p3Phrase := P3Phrase{
			Hash: hash,
			Kind: kind,
			Hint: hint,
		}
		p3Phrases = append(p3Phrases, p3Phrase)
		// ~/~ end
	}
	// ~/~ end

	// ~/~ begin <<docs/index.md#wait-for-ctrl-C>>[init]
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go application(p3Phrases)
	<-sigChan
	fmt.Println()
	log.Println("goodbye, thanks for playing!")
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
		kind = sha256Hk
	default:
		err := fmt.Errorf("could not determine hashkind of '%s'", before)
		return nokind, before, err
	}
	return kind, after, nil
}

// ~/~ end
// ~/~ begin <<docs/index.md#helpers>>[1]
func application(p3Phrases []P3Phrase) {
	for {
		// ~/~ begin <<docs/index.md#shuffle-p3-phrases>>[init]
		rand.Shuffle(len(p3Phrases), func(i, j int) {
			p3Phrases[i], p3Phrases[j] = p3Phrases[j], p3Phrases[i]
		})
		// ~/~ end

		promptP3Phrases(p3Phrases)
	}
}

// ~/~ end
// ~/~ begin <<docs/index.md#helpers>>[2]
func promptP3Phrases(p3Phrases []P3Phrase) {
	for _, p3Phrase := range p3Phrases {
		fmt.Println(p3Phrase.Hint)
		b, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Printf("error on prompting for p3 phrase: %v", err)
			continue
		}

		// ~/~ begin <<docs/index.md#verify-password-and-provide-feedback>>[init]
		foundHashB := sha256.Sum256(b)
		foundHash := fmt.Sprintf("%x", foundHashB)

		if foundHash == p3Phrase.Hash {
			fmt.Println("correct!")
		} else {
			fmt.Println("incorrect!")
		}
		// ~/~ end
	}
}

// ~/~ end
// ~/~ end

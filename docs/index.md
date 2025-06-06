<style>
/* make the code's class, id, and file path bold */
span.filename {
    font-weight: bold;
}
</style>

# Practice Password Phrases

## Requirements

This is **p3** - practice password phrases.
I often mistype passwords, or forget them, so it's a good idea to practice them.
The overarching idea is to store the passwords checksums in a file along with a prompt of sorts.
A password hint is a good prompt.
The program would provide the prompt, the user would enter the password in a no-echo terminal,
and then the program would verify if user got the password correct.

Example

```bash
$ ./p3 --config config.conf
hint: luks password

correct!
hint: root password

incorrect!
hint: luks password
bye, and thanks for playing
```

Note that the blank lines are where the user entered the password.
Since it's a no-echo terminal, no characters are shown.

## Outline

Seems small enough for a single file.

```{.go file=src/main.go}
package main

import (
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/signal"
	"strings"
	"syscall"

	<<hocon-imports>>
)

<<hash-kinds>>

<<define-p3-phrases-type>>

func main() {
	log.Println("starting p3")

	<<setup-flags>>
	flag.Parse()
	<<validate-flags>>

	<<hocon-parse>>

	<<application-loop>>
}

<<helpers>>
```

## Configuration

### Config File

Using [HOCON](https://github.com/lightbend/config/blob/main/HOCON.md) for the configuration.
An example configuration is something like

```{.hocon #config-example file=src/config.conf.sample}
phrases: [
  {
    hint: one,
    hash: sha256-7692c3ad3540bb803c020b3aee66cd8887123234ea0c6e7143c0add73ff431ed,
  },
  {
    hint: two,
    hash: sha256-3fc4ccfe745870e2c0d99f71f30ff0656c8dedd41cc1d7d3d376b0dbe685e2f3,
  }
]
```

In the example above, the passwords are the same as their hints.

The root level must have `phrases` which is an array.
Each array element is an object with two keys, `hint` and `hash`.
`hint` is a string.
`hash` is a string that begins with a recognized hash followed by a `-` followed by a checksum.
So a `sha256` hash would look like `sha256-theActualHash`.

Can find this example in **src/config.conf.sample**.

### Flags

The hocon config file is passed using the `--config` flag

```{.go #setup-flags}
configPath := flag.String("config", "", "path to config file")
```

```{.go #validate-flags}
if len(*configPath) == 0 {
	log.Fatalln("--config flag was provided but no value was given")
}
<<validate-config-exists>>
log.Printf("config file '%s' was provided\n", *configPath)
```

```{.go #validate-config-exists}
if _, err := os.Stat(*configPath); err != nil {
	if os.IsNotExist(err) {
		log.Fatalf("config file '%s' does not exist: %v", *configPath, err)
	} else {
		log.Fatalf("could not stat '%s': %v", *configPath, err)
	}
}
```

### Parse Config File

We'll use the [hocon](https://pkg.go.dev/github.com/gurkankaymak/hocon) package to parse the hocon file.

```{.go #hocon-imports}
"github.com/gurkankaymak/hocon"
"golang.org/x/term"
```

```{.go #hocon-parse}
config, err := hocon.ParseResource(*configPath)
if err != nil {
	log.Fatalf("could not parse hocon file: %v", err)
}
phrases := config.GetArray("phrases")
if phrases == nil {
	log.Fatalln("malformed config file, cannot find 'phrases'")
}

<<initialize-p3-phrases>>

for _, phrase := range phrases {
	configPhrase, err := hocon.ParseString(phrase.String())
	if err != nil {
		log.Fatalf("could not parse phrase element: %v", err)
	}

	<<parse-config-phrase>>

	kind, hash, err := extractHashKind(kindWithHash)
	if err != nil {
		log.Fatalf("could not extract hash in '%s': %v", kindWithHash, err)
	}

	<<add-to-p3-phrases>>
}
```

Parsing `hint` and `hash` is a little tricky because the hocon module surrounds values with `"` at will.
Not entirely certain why it does that since Go's strings are quite powerfull.
Maybe a bug? Anyway,

```{.go #parse-config-phrase}
hintDirty := configPhrase.Get("hint").String()
hint := strings.Trim(hintDirty, "\"")

kindWithHashDirty := configPhrase.Get("hash").String()
kindWithHash := strings.Trim(kindWithHashDirty, "\"")
```

Implementing `extractHashKind` just requires spliting across the first `-`,
verifying that the kind of hash is recognized,
translating the hash kind to the correct enum,
and then returning both the hash kind and the hash value.

First let's define an enum of hash kinds

```{.go #hash-kinds}
type HashKind int

<<hash-kinds-string>>

const (
	nokind   HashKind = -1
	sha256Hk HashKind = iota
)

```

Add a `String()` method that would convert the `HashKind` checksum to a sensible string.
Example output looks like

```
sha256
nokind (error)
```

```{.go #hash-kinds-string}
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

```

Notice the `panic`. That should never get triggered except during development and
someone forgets to extend this method so the `panic` is triggered.

Next let's define the `extractHashKind` function

```{.go #helpers}
func extractHashKind(kindWithHash string) (HashKind, string, error) {
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

```

### Transform to P3 Phrases

Let's use an intermediate data structure to capture the parsed values.
This intermediate data structure can then be passed to the part of the application that tests for the phrases.

```{.go #define-p3-phrases-type}
type P3Phrase struct {
	Hash string
	Hint string
	Kind HashKind
}

<<p3-phrase-string-method>>
```

We'll define a P3Phrase String method that looks like

```
hash=aaaa,hint=one,kind=sha256
```

```{.go #p3-phrase-string-method}
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

```

We'll use the above type to store each entry of `phrases`.

```{.go #initialize-p3-phrases}
p3Phrases := make([]P3Phrase, 0, len(phrases))
```

```{.go #add-to-p3-phrases}
p3Phrase := P3Phrase{
	Hash: hash,
	Kind: kind,
	Hint: hint,
}
p3Phrases = append(p3Phrases, p3Phrase)
```

## Application

### Overview

The application does the following:

```
1. shuffle the list of p3phrases
2. for each p3phrase
   2.1. prompt the hint
   2.2. get input from user (with no-echo)
   2.3. verify and provide feedback
```

Additionally if the user hits cntl-C, the application exits gracefully with a nice message.

### Control C Handling

The main goroutine will wait on ctrl+C to trigger (i.e. `sigint`).
Meanwhile another go routine will actually perform the application.

```{.go #application-loop}
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt)
go application(p3Phrases)
<-sigChan
fmt.Println()
log.Println("goodbye, thanks for playing!")
```

### P3 Phrases Application

```{.go #helpers}
func application(p3Phrases []P3Phrase) {
	for {
		<<shuffle-p3-phrases>>

		promptP3Phrases(p3Phrases)
	}
}

```

```{.go #shuffle-p3-phrases}
rand.Shuffle(len(p3Phrases), func(i, j int) {
	p3Phrases[i], p3Phrases[j] = p3Phrases[j], p3Phrases[i]
})
```

The application itself would use the [term](https://pkg.go.dev/golang.org/x/term) package.
Conveniently, it has the function `term.ReadPassword` which reads a line with no-echo from the terminal.

```{.go #helpers}
func promptP3Phrases(p3Phrases []P3Phrase) {
	for _, p3Phrase := range p3Phrases {
		fmt.Println(p3Phrase.Hint)
		b, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			log.Printf("error on prompting for p3 phrase: %v", err)
			continue
		}

		<<verify-password-and-provide-feedback>>
	}
}

```

```{.go #verify-password-and-provide-feedback}
foundHash, err := computeHash(p3Phrase.Kind, b)
if err != nil {
	log.Printf("error on computing hash: %v", err)
	continue
}

if foundHash == p3Phrase.Hash {
	fmt.Println("correct!")
} else {
	fmt.Println("incorrect!")
}
```

The `computeHash` helper determines the appropriate hash function to use
and returns a string of the hash.

- `sha256Hk` the hash is the sha256 in hexadecimal.

```{.go #helpers}
func computeHash(kind HashKind, b []byte) (string, error) {
	switch kind {
	case sha256Hk:
		foundHashB := sha256.Sum256(b)
		foundHash := fmt.Sprintf("%x", foundHashB)
		return foundHash, nil
	case nokind:
		err := errors.New("cannot compute hash of no kind")
		return "", err
	default:
		s := fmt.Sprintf("cannot compute hash for hash-kind=%d", kind)
		panic(s)
	}
}

```

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

```{.go #main file=src/main.go}
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

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

	<<wait-for-ctrl-C>>
}

<<helpers>>
```

## Configuration

### Config File

Using HOCON for the configuration.
An example configuration is something like

```{.hocon #config-example file=src/config.conf.sample}
phrases: [
  {
    hint: one,
    hash: sha256-aaa,
  },
  {
    hint: two,
    hash: sha256-bbb,
  }
]
```

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

	kind, hash, err := extractHash(kindWithHash)
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

Parsing `extractHash` just requires spliting across the first `-`,
verifying that the kind of hash is recognized,
translating the hash kind to the correct enum,
and then returning both the hash kind and the hash value.

First let's define an enum of hash kinds

```{.go #hash-kinds}
type HashKind int

<<hash-kinds-string>>

const (
	nokind HashKind = -1
	sha256 HashKind = iota
)

```

Add a `String()` method that would convert the `HashKind` checksum to a sensible string.
Example output looks like

```
sha256
nokind (error)
```

Notice the `panic`. That should never get triggered except during development and
someone forgets to extend this method so the `panic` is triggered.

```{.go #hash-kinds-string}
func (hk HashKind) String() string {
	switch hk {
	case sha256:
		return "sha256"
	case nokind:
		return "nokind (error)"
	default:
		s := fmt.Sprintf("cannot convert to string, unexpected HashKind %d", hk)
		panic(s)
	}
}

```

Next let's define the `extractHash` function

```{.go #helpers}
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

1. shuffle the list of p3phrases
1. for each p3phrase
   1. ask user
   1. get input from user
   1. verify and provide feedback

Additionally if the user hits cntl-C, the application exits gracefully with a nice message.

### Control C Handling

The main goroutine will wait on ctrl+C to trigger.
Meanwhile another go routine will actually perform the application.

```{.go #wait-for-ctrl-C}
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
    //TODO
    log.Println(p3Phrases)
}
```

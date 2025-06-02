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
    "log"
)

func main() {
    log.Println("starting p3")
}
```

## Configuration

### Config File

Using HOCON for the configuration.
An example configuration is something like

```{.hocon #config-example file=src/config.conf.sample}
phrases: [
  {
    "hint": "one",
    "hash": "sha256-aaa",
  },
  {
    "hint": "two",
    "hash": "sha256-bbb",
  }
]
```

The root level must have `phrases` which is an array.
Each array element is an object with two keys, `hint` and `hash`.
`hint` is a string.
`hash` is a string that begins with a recognized hash followed by a `-` followed by a checksum.
So a `sha256` hash would look like `sha256-theActualHash`.

Can find this example in **src/config.conf.sample**.

### Parsing the config file

Let's parse the above config file and store it in a data structure.

```{.go #phrases-ds}
type Phrase struct {
    hint String
    hash String
}

```

Each object in the config's `phrases` will populate one `Phrase` instance.

To perform the parsing, use the [hocon module](github.com/gurkankaymak/hocon).

```{.go #parse-imports}
"github.com/gurkankaymak/hocon"
```

Get the path to the configuration file using a flag, say `--config <string>`.
Can simply error out if the flag is missing.

```{.go #std-imports}
"flag"
```

```{.go #configure-flags}
config := flag.String("config", "", "path to config file")
```

Presuming that the configuration gets parsed, the next step is to validate
the config.

```{.go #validate-config}
if config == nil {
    log.Fatalf("no config path was given")
}
```

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
	"log"
	"os"
)

func main() {
	log.Println("starting p3")

	<<setup-flags>>
	flag.Parse()
	<<validate-flags>>
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
```

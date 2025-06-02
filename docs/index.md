<style>
/* make the code's class, id, and file path bold */
span.filename {
    font-weight: bold;
}
</style>

# Practice Password Phrases

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

## Configuration

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

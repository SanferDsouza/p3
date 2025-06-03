# P3

Practice password phrases (p3).
Store the password hashes in a file and practice the passwords.
The program is written using [entangled](https://github.com/entangled/entangled) for literate programming.

If you have `nix` installed, run `nix develop`, followed by `mkdocs serve` to read the code documentation.
Alternatively you can run `virtualenv -p312 venv` followed by `source ./venv/bin/activate`
and `pip install -r requirements.txt` and then `mkdocs serve`.

To run the program, `cd src` and run

```bash
go run main.go --config <path-to-config>
```

Here's a sample config

```
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

There is a lot to be desired in this program.
Nevertheless it is well-documented and less complex compared to without the documentation.
I think I'll use literate programming in my next project!

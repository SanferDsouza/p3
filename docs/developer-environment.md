# Developer Environment

## Tmux

I typically have a tmux window open with three separate panes.
In the first pane I launch `mkdocs serve` which launches the documentation server.
In the second pane I launch `entangled watch` which tangles/stitches docs and code.
Can accomplish this by writing a **.tmux.conf** file and then running

```
tmux new-session -c . -f ~/.tmux.conf \; source-file .tmux.conf
```

which would first use the user's **~/.tmux.conf**
and then would run the command in the local **.tmux.conf**.

```{.tmux file=.tmux.conf}
<<tmux-contents>>
```

First create a new session

```{.tmux #tmux-contents}
new
```

Then we want to run `mkdocs serve` in first (and currently only) pane

```{.tmux #tmux-contents}
send 'mkdocs serve' Enter
```

After that we want to add another pane via a vertical split
and launch `entangled watch` in the new pane

```{.tmux #tmux-contents}
splitw -h
send 'entangled watch' Enter
```

And then we want to again add another pane via a vertical split
and launch `jj st` in this yet another new pane

```{.tmux #tmux-contents}
splitw -h
send 'jj st' Enter
```

Finally we want to evenly space the vertical panes

```{.tmux #tmux-contents}
selectl even-horizontal
```

A side note, some folks prefer `C-m` instead of `Enter`.
Imo `Enter` is more readable.

## Nix Flake

Using **Nix** is a convenient way to ensure a reproducible environment.
I use nix flakes because the versions are *locked* and resolution is cached (so sub-second environment activation).
Outline:

```{.nix file=flake.nix}
{
  description = ''
    <<flake-description>>
  '';

  inputs = {
    <<flake-inputs>>
  };

  outputs =
    {
      <<flake-output-args>>
      self,
    }:
    <<flake-output-body>>
    
}
```

```{.nix #flake-description}
Practice typing passwords by comparing them to their hashes.
Can practice multiple passwords by providing a description
that identifies each password hash. Can be a hint.
Never type a password incorrectly again (after some practice)!
```

To support multiple architectures (although I personally use `x86_64-linux`),
we'll need `flake-utils`

```{.nix #flake-inputs}
flake-utils.url = "github:numtide/flake-utils";
```

```{.nix #flake-output-args}
flake-utils,
```

```{.nix #flake-output-body}
flake-utils.lib.eachDefaultSystem (
  system:
  <<flake-utils-body>>
);

```

In the body of the `flake-utils` we'll want to import `nixpkgs`
and then use it to load packages to the devshell.

```{.nix #flake-inputs}
nixpkgs.url = "github:nixos/nixpkgs/24.11";
```

```{.nix #flake-output-args}
nixpkgs,
```

```{.nix #flake-utils-body}
let
  pkgs = import nixpkgs { inherit system; };
in
{
  devShells.default = pkgs.mkShellNoCC {
    packages = with pkgs; [
      <<dev-shells-pkgs>>
    ];
    shellHook = ''
      <<dev-shells-shell-hook>>
    '';
  };
}
```

We'll need the following packages:

- `go` since the source code is in go
- `python312` since entangled and mkdocs needs at least version 312
- `tmux` which is a developer utility. Since there's a [tmux section](#tmux), might as well include it for completeness
- `virtualenv` since that would create the virtual environment

```{.nix #dev-shells-pkgs}
go
python312
tmux
virtualenv
```

Then we'll setup the shell hook to create the virtualenv and install the contents of **requirements.txt** into it.

```{.nix #dev-shells-shell-hook}
virtualenv -p312 venv
source ./venv/bin/activate
pip install -qr requirements.txt
```

To activate the environment, run `nix develop`.

## Direnv

I use [direnv](https://direnv.net/) to call `nix develop` every time I enter the current directory.
For that to work as desired, need a **.envrc** file.

```{.bash file=.envrc}
watch_file flake.nix
watch_file flake.lock
eval "$(nix print-dev-env)"
```

See also [direnv with nix](https://github.com/direnv/direnv/wiki/Nix#hand-rolled-nix-flakes-integration).

Would need to run `direnv allow` while inside the repo the first time around.

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

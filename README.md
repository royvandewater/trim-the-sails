# trim-the-sails

Iterates over every bare git repo in the current directory (directories whose
names end in `.git`) and prunes their branches:

1. `git fetch --all --prune` — drops remote-tracking refs for branches deleted
   on the remote.
2. Deletes any local branch whose upstream is now gone.

## Usage

```bash
cd /path/to/dir/of/bare/repos
trim-the-sails
```

## Build

```bash
go build -o trim-the-sails .
```

## Test

```bash
go test ./...
```

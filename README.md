# gitcsv

Store git log as CSV

## Install

```
go get github.com/mcls/gitcsv
```

## Usage

```
gitcsv path/to/git/repo/ ./git_log.csv
```

## Example

Parsing the [linux](https://github.com/torvalds/linux) kernel's git history:

```
$ gitcsv ~/Desktop/experiment/linux
2017/03/06 09:31:26 Opening repo
2017/03/06 09:31:26 Parsing each commit in history
2017/03/06 09:31:52 Done.
```

As you can see it takes less than 30 seconds to parse more than 600k commits.

## Dependencies

This uses [git2go](https://github.com/libgit2/git2go), so it depends on [libgit2](https://libgit2.github.com/).

To install with Homebrew run `brew install libgit2`.

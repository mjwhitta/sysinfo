# sysinfo

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?style=for-the-badge&logo=cookiecutter)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/sysinfo)](https://goreportcard.com/report/github.com/mjwhitta/sysinfo)

## What is this?

Provides system info for use by other tools.

## How to install

Open a terminal and run the following:

```
$ go get --ldflags "-s -w" --trimpath -u github.com/mjwhitta/sysinfo
$ go install --ldflags "-s -w" --trimpath \
    github.com/mjwhitta/sysinfo/cmd/sysinfo@latest
```

Or compile from source:

```
$ git clone https://github.com/mjwhitta/sysinfo.git
$ cd sysinfo
$ git submodule update --init
$ make
```

## How to use

```
$ sysinfo
```

```
package main

import (
    "fmt"

    "github.com/mjwhitta/sysinfo"
)

func main() {
    fmt.Println(sysinfo.New())
}
```

## Configuration

Configuration is stored in `$HOME/.config/sysinfo/rc`. The default
config looks like:

```
{
  "dataColors": [
    "green"
  ],
  "fieldColors": [
    "blue"
  ]
}
```

These values can be adjusted to meet your needs.

## Links

- [Source](https://github.com/mjwhitta/sysinfo)

## TODO

- Better README.md

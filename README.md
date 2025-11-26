# sysinfo

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/sysinfo?style=for-the-badge)](https://goreportcard.com/report/github.com/mjwhitta/sysinfo)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/mjwhitta/sysinfo/ci.yaml?style=for-the-badge)](https://github.com/mjwhitta/sysinfo/actions)
![License](https://img.shields.io/github/license/mjwhitta/sysinfo?style=for-the-badge)

## What is this?

Provides system info for use by other tools.

## How to install

Open a terminal and run the following:

```
$ go get -u github.com/mjwhitta/sysinfo
$ go install github.com/mjwhitta/sysinfo/cmd/sysinfo@latest
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
- Windows performance improvements
    - Use API, not powershell

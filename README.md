# sysinfo

Provides system info for use by other tools.

## How to install

Open a terminal and run the following:

```
$ go get -u gitlab.com/mjwhitta/sysinfo/cmd/sysinfo
```

Or install from source:

```
$ git clone https://gitlab.com/mjwhitta/sysinfo.git
$ cd sysinfo
$ make
$ cp ./build/sysinfo ~/.local/bin/
```

## How to use

```
$ sysinfo
```

```
package main

import (
    "fmt"

    "gitlab.com/mjwhitta/sysinfo"
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
  "kbg": "on_default",
  "kfg": "blue",
  "vbg": "on_default",
  "vfg": "green"
}
```

These values can be changed to adjust the key/value bg/fg colors.

## Links

- [Source](https://gitlab.com/mjwhitta/sysinfo)

## TODO

- Better README.md

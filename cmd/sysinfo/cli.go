package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mjwhitta/cli"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/sysinfo"
)

// Exit status
const (
	Good = iota
	InvalidOption
	MissingOption
	InvalidArgument
	MissingArgument
	ExtraArgument
	Exception
)

// Flags
var flags struct {
	fields  cli.StringList
	nocolor bool
	verbose bool
	version bool
}

func init() {
	// Configure cli package
	cli.Align = true
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = filepath.Base(os.Args[0]) + " [OPTIONS]"
	cli.BugEmail = "sysinfo.bugs@whitta.dev"

	cli.ExitStatus(
		"Normally the exit status is 0. In the event of an error the",
		"exit status will be one of the below:\n\n",
		fmt.Sprintf("%d: Invalid option\n", InvalidOption),
		fmt.Sprintf("%d: Missing option\n", MissingOption),
		fmt.Sprintf("%d: Invalid argument\n", InvalidArgument),
		fmt.Sprintf("%d: Missing argument\n", MissingArgument),
		fmt.Sprintf("%d: Extra argument\n", ExtraArgument),
		fmt.Sprintf("%d: Exception", Exception),
	)
	cli.Info(
		"System information at a glance. Configuration is stored in",
		"~/.config/sysinfo/rc.",
	)
	cli.SectionAligned(
		"FIELDS",
		":",
		"blank:Blank line\n",
		"colors:Sample of terminal colors\n",
		"cpu:CPU info\n",
		"fs:Filesystem usage\n",
		"host:Hostname\n",
		"ip:IPv4/IPv6 addresses\n",
		"kernel:Kernel info\n",
		"os:Operating System info\n",
		"ram:RAM usage\n",
		"shell:Current shell\n",
		"tty:TTY info\n",
		"uptime:Uptime",
	)

	cli.Title = "SysInfo"

	// Parse cli flags
	cli.Flag(
		&flags.fields,
		"f",
		"field",
		"Show specified field. Can be used more than once. By",
		"default, all fields are shown. Use this flag to adjust the",
		"order.",
	)
	cli.Flag(
		&flags.nocolor,
		"no-color",
		false,
		"Disable colorized output.",
	)
	cli.Flag(
		&flags.verbose,
		"v",
		"verbose",
		false,
		"Show stacktrace, if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	hl.Disable(flags.nocolor)

	// Short circuit if version was requested
	if flags.version {
		fmt.Println(
			filepath.Base(os.Args[0]) + " version " + sysinfo.Version,
		)
		os.Exit(Good)
	}

	// Validate cli flags
	if cli.NArg() > 1 {
		cli.Usage(ExtraArgument)
	}
}

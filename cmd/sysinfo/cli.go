package main

import (
	"os"
	"strings"

	"gitlab.com/mjwhitta/cli"
	hl "gitlab.com/mjwhitta/hilighter"
)

// Flags
type cliFlags struct {
	nocolor bool
	verbose bool
	version bool
}

var flags cliFlags

func init() {
	// Configure cli package
	cli.Align = true
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = hl.Sprintf("%s [OPTIONS]", os.Args[0])
	cli.BugEmail = "sysinfo.bugs@whitta.dev"
	cli.ExitStatus = strings.Join(
		[]string{
			"Normally the exit status is 0. In the event of an error",
			"the exit status will be one of the below:\n\n",
			"1: Invalid option\n",
			"2: Invalid argument\n",
			"3: Extra arguments\n",
			"4: Exception",
		},
		" ",
	)
	cli.Info = "System information at a glance."
	cli.Title = "SysInfo"

	// Parse cli flags
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
		"Show show stacktrace if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	// Validate cli flags
	if cli.NArg() > 1 {
		cli.Usage(ExtraArguments)
	}
}

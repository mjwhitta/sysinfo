package main

import (
	"os"
	"strings"

	"gitlab.com/mjwhitta/cli"
	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/sysinfo"
)

// Flags
type cliFlags struct {
	fields  cli.StringList
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
	cli.Section(
		"FIELDS",
		strings.Join(
			[]string{
				"blank: Use a blank line as a separator\n",
				"colors: Show terminal colors\n",
				"cpu: Show cpu info\n",
				"fs: Show filesystem usage\n",
				"host: Show hostname\n",
				"ipv4: Show IPv4 addresses\n",
				"ipv6: Show IPv6 addresses\n",
				"kernel: Show kernel info\n",
				"os: Show operating system info\n",
				"ram: Show RAM usage\n",
				"shell: Show current shell\n",
				"tty: Show TTY info\n",
				"uptime: Show uptime",
			},
			"",
		),
	)
	cli.Info = "System information at a glance."
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
		"Show show stacktrace if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	hl.Disable(flags.nocolor)

	// Short circuit if version was requested
	if flags.version {
		hl.Printf("sysinfo version %s\n", sysinfo.Version)
		os.Exit(Good)
	}

	// Validate cli flags
	if cli.NArg() > 1 {
		cli.Usage(ExtraArguments)
	}
}

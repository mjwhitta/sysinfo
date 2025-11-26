package sysinfo

import "regexp"

// Version is the package version
const Version string = "1.7.1"

var (
	reCommaSpace *regexp.Regexp = regexp.MustCompile(
		`^(,?\s*)+|(,?\s*)+$`,
	)
	reCPUBrand *regexp.Regexp = regexp.MustCompile(
		`\((R|TM)\)| (@|CPU)`,
	)
	reHrMin *regexp.Regexp = regexp.MustCompile(
		`0?(\d+):0?(\d+)`,
	)
	reModelName *regexp.Regexp = regexp.MustCompile(
		`(cpu model|model name)\s+:\s+(.+)`,
	)
	reOneHrMin *regexp.Regexp = regexp.MustCompile(
		`(^|,\s+)1\s+(hour|min)s`,
	)
	rePrettyName *regexp.Regexp = regexp.MustCompile(
		`PRETTY_NAME="(.+)"`,
	)
	reRAM *regexp.Regexp = regexp.MustCompile(
		`Mem:\s+(\d+)\s+(\d+)`,
	)
	reUptimeEnds *regexp.Regexp = regexp.MustCompile(
		`^.*up\s+|,\s+\d+\s+user.+$`,
	)
	reWhiteSpace  *regexp.Regexp = regexp.MustCompile(`\s+`)
	reZeroHrsMins *regexp.Regexp = regexp.MustCompile(
		`(^|,\s+)0\s+(hour|min)s`,
	)
	titleCase map[string]string = map[string]string{
		"cpu":    "CPU",
		"homefs": "HomeFS",
		"host":   "Host",
		"ipv4":   "IPv4",
		"ipv6":   "IPv6",
		"kernel": "Kernel",
		"os":     "OS",
		"ram":    "RAM",
		"rootfs": "RootFS",
		"shell":  "Shell",
		"tty":    "TTY",
		"uptime": "Uptime",
	}
)

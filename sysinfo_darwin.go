//go:build darwin

package sysinfo

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	hl "github.com/mjwhitta/hilighter"
)

func (s *SysInfo) colors() {
	s.Colors = strings.Join(
		[]string{
			hl.Hilights([]string{"light_black", "on_black"}, "▄▄▄"),
			hl.Hilights([]string{"light_red", "on_red"}, "▄▄▄"),
			hl.Hilights([]string{"light_green", "on_green"}, "▄▄▄"),
			hl.Hilights([]string{"light_yellow", "on_yellow"}, "▄▄▄"),
			hl.Hilights([]string{"light_blue", "on_blue"}, "▄▄▄"),
			hl.Hilights(
				[]string{"light_magenta", "on_magenta"},
				"▄▄▄",
			),
			hl.Hilights([]string{"light_cyan", "on_cyan"}, "▄▄▄"),
			hl.Hilights([]string{"light_white", "on_white"}, "▄▄▄"),
		},
		"",
	)
}

func (s *SysInfo) cpu() {
	s.CPU = s.exec("sysctl", "-n", "machdep.cpu.brand_string")
	s.CPU = reCPUBrand.ReplaceAllString(s.CPU, "")
	s.CPU = reWhiteSpace.ReplaceAllString(s.CPU, " ")
}

func (s *SysInfo) filesystems() {
	s.RootFS = s.fsUsage("/")

	if s.HomeFS = s.fsUsage("/home"); s.HomeFS == s.RootFS {
		s.HomeFS = ""
	}

	if s.RootFS == "" {
		s.RootFS = "unknown"
	}
}

func (s *SysInfo) fsUsage(path string) string {
	var cols []string
	var usage string = s.exec("df", "-h", path)

	for _, line := range strings.Split(usage, "\n") {
		cols = strings.Fields(line)

		//nolint:mnd // Validate output format
		if (len(cols) == 9) && (cols[8] == path) {
			return cols[2] + " / " + cols[1] + " (" + cols[4] + ")"
		}
	}

	return ""
}

func (s *SysInfo) kernel() {
	s.Kernel = s.exec("sysctl", "-n", "kern.osrelease")
}

func (s *SysInfo) operatingSystem() {
	s.OS = s.exec("uname", "-m", "-s")
}

func (s *SysInfo) ram() {
	var e error
	var mb int = 1024 * 1024
	var phys int
	var tmp string
	var total int
	var user int

	s.RAM = "unknown"

	tmp = s.exec("sysctl", "-n", "hw.physmem")
	if phys, e = strconv.Atoi(tmp); e != nil {
		return
	}

	tmp = s.exec("sysctl", "-n", "hw.usermem")
	if user, e = strconv.Atoi(tmp); e != nil {
		return
	}

	tmp = s.exec("sysctl", "-n", "hw.memsize")
	if total, e = strconv.Atoi(tmp); e != nil {
		return
	}

	s.RAM = fmt.Sprintf("%d MB / %d MB", (phys+user)/mb, total/mb)
}

func (s *SysInfo) shell() {
	s.Shell = "unknown"

	if sh, ok := os.LookupEnv("SHELL"); ok {
		s.Shell = strings.TrimSpace(sh)
	}
}

func (s *SysInfo) tty() {
	// There's probably a better way
	s.TTY = os.Getenv("GPG_TTY")
	if s.TTY = strings.TrimSpace(s.TTY); s.TTY == "" {
		s.TTY = "unknown"
	}
}

func (s *SysInfo) uptime() {
	var uptime string

	s.Uptime = "0 mins"

	// Fail fast
	if uptime = s.exec("uptime"); uptime == "" {
		return
	}

	// Strip extra whitespace
	uptime = reWhiteSpace.ReplaceAllString(uptime, " ")

	// Strip leading and trailing data
	uptime = reUptimeEnds.ReplaceAllString(uptime, "")

	// Make plural, if not already
	if strings.HasSuffix(uptime, "min") {
		uptime += "s"
	}

	// Convert hours:mins to match days
	uptime = reHrMin.ReplaceAllString(uptime, "$1 hours, $2 mins")

	// Remove 0 hours and mins
	uptime = reZeroHrsMins.ReplaceAllString(uptime, "")

	// Make singular, if needed
	uptime = reOneHrMin.ReplaceAllString(uptime, "${1}1 ${2}")

	// Remove leading and trailing commas or whitespace
	uptime = reCommaSpace.ReplaceAllString(uptime, "")

	if uptime != "" {
		s.Uptime = uptime
	}
}

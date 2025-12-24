//go:build !darwin && !windows

package sysinfo

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/pathname"
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
	var e error
	var info []byte
	var m [][]string

	s.CPU = "unknown"

	if info, e = os.ReadFile("/proc/cpuinfo"); e != nil {
		return
	}

	m = reModelName.FindAllStringSubmatch(string(info), -1)
	if len(m) > 0 {
		s.CPU = fmt.Sprintf(
			"%s(x%d)",
			reCPUBrand.ReplaceAllString(m[0][2], ""),
			len(m),
		)
		s.CPU = reWhiteSpace.ReplaceAllString(s.CPU, " ")
	}
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
		if (len(cols) == 6) && (cols[5] == path) {
			return cols[2] + " / " + cols[1] + " (" + cols[4] + ")"
		}
	}

	return ""
}

func (s *SysInfo) kernel() {
	var b []byte

	s.Kernel = "unknown"

	b, _ = os.ReadFile("/proc/sys/kernel/osrelease")
	if b = bytes.TrimSpace(b); len(b) > 0 {
		s.Kernel = string(b)
	}
}

func (s *SysInfo) operatingSystem() {
	var b []byte
	var e error
	var m [][]string

	s.OS = s.exec("uname", "-m", "-s")

	if ok, _ := pathname.DoesExist("/etc/os-release"); ok {
		if b, e = os.ReadFile("/etc/os-release"); e != nil {
			return
		}

		m = rePrettyName.FindAllStringSubmatch(string(b), -1)
		if len(m) > 0 {
			s.OS = m[0][1] + " " + s.exec("uname", "-m")
		}
	}
}

func (s *SysInfo) ram() {
	var m [][]string
	var mb int = 1024
	var total int
	var used int

	s.RAM = "unknown"

	m = reRAM.FindAllStringSubmatch(s.exec("free"), -1)
	if len(m) > 0 {
		// No need to check the errors here b/c the regex capture
		// group has to be an int
		total, _ = strconv.Atoi(m[0][1])
		used, _ = strconv.Atoi(m[0][2])

		s.RAM = fmt.Sprintf("%d MB / %d MB", used/mb, total/mb)
	}
}

func (s *SysInfo) shell() {
	s.Shell = "unknown"

	if sh, ok := os.LookupEnv("SHELL"); ok {
		s.Shell = strings.TrimSpace(sh)
	}
}

func (s *SysInfo) tty() {
	s.TTY = "unknown"
	if tty, e := os.Readlink("/proc/self/fd/0"); e == nil {
		s.TTY = strings.TrimSpace(tty)
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

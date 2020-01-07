package sysinfo

import (
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/where"
)

// SysInfo is a struct contained relevant system information.
type SysInfo struct {
	Colors   string
	CPU      string
	HomeFS   string
	Hostname string
	IPv4     string
	IPv6     string
	Kernel   string
	order    []string
	OS       string
	RAM      string
	RootFS   string
	Shell    string
	TTY      string
	Uptime   string
}

// New will return a SysInfo pointer. A list of fields can be
// supplied if all info is not wanted.
func New(fields ...string) *SysInfo {
	var s = &SysInfo{}

	s.order = fields
	if len(fields) == 0 {
		s.order = []string{
			"hostname",
			"os",
			"kernel",
			"uptime",
			"ip",
			"shell",
			"tty",
			"cpu",
			"ram",
			"fs",
			"colors",
		}
	}

	for _, field := range s.order {
		switch field {
		case "colors":
			s.colors()
		case "cpu":
			s.cpu()
		case "fs":
			s.filesystems()
		case "host", "hostname":
			s.hostname()
		case "ip":
			s.ipv4()
			s.ipv6()
		case "ipv4":
			s.ipv4()
		case "ipv6":
			s.ipv6()
		case "kernel":
			s.kernel()
		case "os":
			s.operatingSystem()
		case "ram":
			s.ram()
		case "shell":
			s.shell()
		case "tty":
			s.tty()
		case "uptime":
			s.uptime()
		}
	}

	return s
}

func (s *SysInfo) exec(cmd string, cli ...string) string {
	var e error
	var o []byte

	if len(cmd) == 0 || len(where.Is(cmd)) == 0 {
		return ""
	}

	if o, e = exec.Command(cmd, cli...).Output(); e != nil {
		// return e.Error()
		panic(e)
	}

	return strings.TrimSpace(string(o))
}

func (s *SysInfo) colors() string {
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
	return s.Colors
}

func (s *SysInfo) cpu() string {
	var e error
	var info []byte
	var r *regexp.Regexp

	if info, e = ioutil.ReadFile("/proc/cpuinfo"); e != nil {
		panic(e)
	}

	r = regexp.MustCompile(`name\s+:\s+(.+)`)
	for _, match := range r.FindAllStringSubmatch(string(info), -1) {
		s.CPU = match[1]
		break
	}

	r = regexp.MustCompile(`\((R|TM)\)| (@|CPU)`)
	s.CPU = r.ReplaceAllString(s.CPU, "")

	return s.CPU
}

func (s *SysInfo) filesystems() []string {
	var out = []string{}

	// TODO filesystems

	if len(s.HomeFS) > 0 {
		out = append(out, s.HomeFS)
	}

	if len(s.RootFS) > 0 {
		out = append(out, s.RootFS)
	}

	return out
}

func (s *SysInfo) hostname() string {
	s.Hostname = s.exec("hostname", "-s")
	return s.Hostname
}

func (s *SysInfo) ipv4() string {
	var dev string
	var matches [][]string
	var r *regexp.Regexp

	r = regexp.MustCompile(`^default.+dev\s+(\S+)`)
	matches = r.FindAllStringSubmatch(s.exec("ip", "r"), -1)
	for _, match := range matches {
		dev = match[1]
		break
	}

	r = regexp.MustCompile(`(?i)\s+inet\s+(\S+)`)
	matches = r.FindAllStringSubmatch(
		s.exec("ip", "-o", "a", "s", dev),
		-1,
	)
	for _, match := range matches {
		s.IPv4 = match[1]
		break
	}

	return s.IPv4
}

func (s *SysInfo) ipv6() string {
	var dev string
	var matches [][]string
	var r *regexp.Regexp

	r = regexp.MustCompile(`^default.+dev\s+(\S+)`)
	matches = r.FindAllStringSubmatch(s.exec("ip", "r"), -1)
	for _, match := range matches {
		dev = match[1]
		break
	}

	r = regexp.MustCompile(`(?i)\s+inet6\s+(\S+)`)
	matches = r.FindAllStringSubmatch(
		s.exec("ip", "-o", "a", "s", dev),
		-1,
	)
	for _, match := range matches {
		s.IPv6 = match[1]
		break
	}

	r = regexp.MustCompile(`(?i)^fe[89ab]`)
	if r.MatchString(s.IPv6) {
		s.IPv6 = ""
	}

	return s.IPv6
}

func (s *SysInfo) kernel() string {
	s.Kernel = s.exec("uname", "-r")
	return s.Kernel
}

func (s *SysInfo) operatingSystem() string {
	// TODO os
	return s.OS
}

func (s *SysInfo) ram() string {
	// TODO ram
	return s.RAM
}

func (s *SysInfo) shell() string {
	var exists bool
	var sh string

	if sh, exists = os.LookupEnv("SHELL"); exists {
		s.Shell = strings.TrimSpace(sh)
	}

	return s.Shell
}

func (s *SysInfo) tty() string {
	var e error
	if s.TTY, e = os.Readlink("/proc/self/fd/0"); e != nil {
		panic(e)
	}
	return s.TTY
}

func (s *SysInfo) uptime() string {
	// TODO uptime
	return s.Uptime
}

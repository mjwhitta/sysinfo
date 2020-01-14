package sysinfo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/pathname"
	"gitlab.com/mjwhitta/where"
)

// SysInfo is a struct containing relevant system information.
type SysInfo struct {
	Colors string   `json:"-"`
	CPU    string   `json:"CPU,omitempty"`
	Height int      `json:"-"`
	HomeFS string   `json:"HomeFS,omitempty"`
	Host   string   `json:"Host,omitempty"`
	IPv4   string   `json:"IPv4,omitempty"`
	IPv6   string   `json:"IPv6,omitempty"`
	Kernel string   `json:"Kernel,omitempty"`
	order  []string `json:"-"`
	OS     string   `json:"OS,omitempty"`
	RAM    string   `json:"RAM,omitempty"`
	RootFS string   `json:"RootFS,omitempty"`
	Shell  string   `json:"Shell,omitempty"`
	TTY    string   `json:"TTY,omitempty"`
	Uptime string   `json:"Uptime,omitempty"`
	Width  int      `json:"-"`
}

// New will return a SysInfo pointer. A list of fields can be
// supplied if all info is not wanted.
func New(fields ...string) *SysInfo {
	var s = &SysInfo{}

	s.order = fields
	if len(fields) == 0 {
		s.order = []string{
			"Host",
			"OS",
			"Kernel",
			"Uptime",
			"IP",
			"Shell",
			"TTY",
			"CPU",
			"RAM",
			"FS",
			"Colors",
		}
	}

	for _, field := range s.order {
		switch field {
		case "Colors":
			s.colors()
		case "CPU":
			s.cpu()
		case "FS":
			s.filesystems()
		case "Host":
			s.hostname()
		case "IP":
			s.ipv4()
			s.ipv6()
		case "IPv4":
			s.ipv4()
		case "IPv6":
			s.ipv6()
		case "Kernel":
			s.kernel()
		case "OS":
			s.operatingSystem()
		case "RAM":
			s.ram()
		case "Shell":
			s.shell()
		case "TTY":
			s.tty()
		case "Uptime":
			s.uptime()
		default:
			panic(errors.New("Invalid field: " + field))
		}
	}

	for _, line := range strings.Split(hl.Plain(s.String()), "\n") {
		s.Height++
		if len([]rune(line)) > s.Width {
			s.Width = len([]rune(line))
		}
	}

	return s
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

	r = regexp.MustCompile(`(cpu model|model name)\s+:\s+(.+)`)
	for _, match := range r.FindAllStringSubmatch(string(info), -1) {
		s.CPU = match[2]
		break
	}

	r = regexp.MustCompile(`\((R|TM)\)| (@|CPU)`)
	s.CPU = r.ReplaceAllString(s.CPU, "")

	return s.CPU
}

func (s *SysInfo) exec(cmd string, cli ...string) string {
	var e error
	var o []byte

	if len(cmd) == 0 || len(where.Is(cmd)) == 0 {
		return ""
	}

	if o, e = exec.Command(cmd, cli...).Output(); e != nil {
		// return e.Error()
		return ""
		// panic(e)
	}

	return strings.TrimSpace(string(o))
}

func (s *SysInfo) filesystems() []string {
	var out = []string{}

	s.RootFS = s.fsUsage("/")
	s.HomeFS = s.fsUsage("/home")

	if len(s.RootFS) > 0 {
		out = append(out, s.RootFS)
	}

	if (len(s.HomeFS) > 0) && (s.HomeFS != s.RootFS) {
		out = append(out, s.HomeFS)
	} else {
		s.HomeFS = ""
	}

	return out
}

func formatLine(k string, v string, max int) string {
	var line string
	var r = regexp.MustCompile(`%`)

	v = r.ReplaceAllString(v, "%%")

	line = " "
	for i := 0; i < max-len(k); i++ {
		line += " "
	}
	line += hl.Blue(k+":") + " " + hl.LightGreen(v)

	return line
}

func (s *SysInfo) fsUsage(path string) string {
	var matches [][]string
	var r *regexp.Regexp
	var usage string

	usage = s.exec("df", "-h", path)

	r = regexp.MustCompile(`/\S+\s+(\S+)\s+(\S+)\s+\S+\s+(\S+)`)
	matches = r.FindAllStringSubmatch(usage, -1)
	for _, match := range matches {
		return match[2] + " / " + match[1] + " (" + match[3] + ")"
		break
	}

	return ""
}

func (s *SysInfo) hostname() string {
	var e error
	var host []byte

	host, e = ioutil.ReadFile("/proc/sys/kernel/hostname")
	if e == nil {
		s.Host = strings.TrimSpace(string(host))
	}

	return s.Host
}

func (s *SysInfo) ipv4() string {
	var dev string
	var matches [][]string
	var r *regexp.Regexp

	r = regexp.MustCompile(`^default.+dev\s+(\S+)`)
	matches = r.FindAllStringSubmatch(s.exec("ip", "-o", "r"), -1)
	for _, match := range matches {
		dev = match[1]
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
	var e error
	var kernel []byte

	kernel, e = ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if e == nil {
		s.Kernel = strings.TrimSpace(string(kernel))
	}

	return s.Kernel
}

func (s *SysInfo) operatingSystem() string {
	var e error
	var matches [][]string
	var r *regexp.Regexp
	var release []byte

	if pathname.DoesExist("/etc/os-release") {
		if release, e = ioutil.ReadFile("/etc/os-release"); e != nil {
			panic(e)
		}

		r = regexp.MustCompile(`PRETTY_NAME="(.+)"`)
		matches = r.FindAllStringSubmatch(string(release), -1)
		for _, match := range matches {
			s.OS = match[1] + " " + s.exec("uname", "-m")
			break
		}
	} else {
		s.OS = s.exec("uname", "-m", "-s")
	}

	return s.OS
}

func (s *SysInfo) ram() string {
	var matches [][]string
	var r *regexp.Regexp
	var total int
	var used int

	r = regexp.MustCompile(`Mem:\s+(\d+)\s+(\d+)`)
	matches = r.FindAllStringSubmatch(s.exec("free"), -1)
	for _, match := range matches {
		total, _ = strconv.Atoi(match[1])
		used, _ = strconv.Atoi(match[2])

		total /= 1024
		used /= 1024

		s.RAM = strconv.Itoa(used) + " MB"
		s.RAM += " / "
		s.RAM += strconv.Itoa(total) + " MB"

		break
	}

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

// String will convert the SysInfo struct to a printable string.
func (s *SysInfo) String() string {
	var data = map[string]string{}
	var e error
	var hasKey bool
	var max int
	var out []string
	var tmp []byte

	if tmp, e = json.Marshal(s); e != nil {
		panic(e)
	}

	if e = json.Unmarshal(tmp, &data); e != nil {
		panic(e)
	}

	for k := range data {
		if len(k) > max {
			max = len(k)
		}
	}

	for _, field := range s.order {
		switch field {
		case "Colors":
			out = append(out, "")
			out = append(out, " "+s.Colors)
		case "FS":
			field = "RootFS"
			if _, hasKey = data[field]; hasKey {
				out = append(out, formatLine(field, data[field], max))
			}

			field = "HomeFS"
			if _, hasKey = data[field]; hasKey {
				out = append(out, formatLine(field, data[field], max))
			}
		case "IP":
			field = "IPv4"
			if _, hasKey = data[field]; hasKey {
				out = append(out, formatLine(field, data[field], max))
			}

			field = "IPv6"
			if _, hasKey = data[field]; hasKey {
				out = append(out, formatLine(field, data[field], max))
			}
		default:
			out = append(out, formatLine(field, data[field], max))
		}
	}

	return strings.Join(out, "\n")
}

func (s *SysInfo) tty() string {
	var e error
	if s.TTY, e = os.Readlink("/proc/self/fd/0"); e != nil {
		panic(e)
	}
	return s.TTY
}

func (s *SysInfo) uptime() string {
	var r *regexp.Regexp

	s.Uptime = s.exec("uptime")

	r = regexp.MustCompile(`^.+up\s+|,\s+\d+\s+user.+$`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "")

	r = regexp.MustCompile(`(days?),\s+`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "$1, ")

	r = regexp.MustCompile(`0?(\d+):0?(\d+)`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "$1 hour, $2 min")

	r = regexp.MustCompile(`(0 hour, |, 0 min)`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "")

	r = regexp.MustCompile(`((\d\d+|[2-9]) (hour|min))`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "${1}s")

	r = regexp.MustCompile(`\s+`)
	s.Uptime = r.ReplaceAllString(s.Uptime, " ")

	return s.Uptime
}

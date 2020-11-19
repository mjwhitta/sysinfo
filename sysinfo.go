package sysinfo

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/pathname"
	"gitlab.com/mjwhitta/where"
)

// SysInfo is a struct containing relevant system information.
type SysInfo struct {
	Colors      string `json:"-"`
	CPU         string `json:"CPU,omitempty"`
	dataColors  []string
	fieldColors []string
	Height      int    `json:"-"`
	HomeFS      string `json:"HomeFS,omitempty"`
	Host        string `json:"Host,omitempty"`
	IPv4        string `json:"IPv4,omitempty"`
	IPv6        string `json:"IPv6,omitempty"`
	Kernel      string `json:"Kernel,omitempty"`
	order       []string
	OS          string `json:"OS,omitempty"`
	RAM         string `json:"RAM,omitempty"`
	RootFS      string `json:"RootFS,omitempty"`
	Shell       string `json:"Shell,omitempty"`
	TTY         string `json:"TTY,omitempty"`
	Uptime      string `json:"Uptime,omitempty"`
	Width       int    `json:"-"`
}

// New will return a SysInfo pointer. A list of fields can be
// supplied if all info is not wanted.
func New(fields ...string) *SysInfo {
	var s = &SysInfo{}

	s.order = fields
	if len(fields) == 0 {
		s.order = []string{
			"host",
			"os",
			"kernel",
			"uptime",
			"ipv4",
			"ipv6",
			"shell",
			"tty",
			"cpu",
			"ram",
			"fs",
			"blank",
			"colors",
		}
	}

	s.Collect()

	return s
}

func (s *SysInfo) calcSize() {
	s.Height = 0
	s.Width = 0

	for _, line := range strings.Split(hl.Plain(s.String()), "\n") {
		s.Height++
		if len([]rune(line)) > s.Width {
			s.Width = len([]rune(line))
		}
	}
}

// Clear will remove all system info.
func (s *SysInfo) Clear() {
	s.Colors = ""
	s.CPU = ""
	s.HomeFS = ""
	s.Host = ""
	s.IPv4 = ""
	s.IPv6 = ""
	s.Kernel = ""
	s.OS = ""
	s.RAM = ""
	s.RootFS = ""
	s.Shell = ""
	s.TTY = ""
	s.Uptime = ""
	s.calcSize()
}

// Collect will get requested system info.
func (s *SysInfo) Collect() {
	var newOrder []string
	var wg = sync.WaitGroup{}

	for _, field := range s.order {
		switch strings.ToLower(field) {
		case "blank":
			newOrder = append(newOrder, "Blank")
		case "colors":
			newOrder = append(newOrder, "Colors")

			wg.Add(1)
			go func() {
				s.colors()
				wg.Done()
			}()
		case "cpu":
			newOrder = append(newOrder, "CPU")

			wg.Add(1)
			go func() {
				s.cpu()
				wg.Done()
			}()
		case "fs":
			newOrder = append(newOrder, "FS")

			wg.Add(1)
			go func() {
				s.filesystems()
				wg.Done()
			}()
		case "host":
			newOrder = append(newOrder, "Host")

			wg.Add(1)
			go func() {
				s.hostname()
				wg.Done()
			}()
		case "ipv4":
			newOrder = append(newOrder, "IPv4")

			wg.Add(1)
			go func() {
				s.ipv4()
				wg.Done()
			}()
		case "ipv6":
			newOrder = append(newOrder, "IPv6")

			wg.Add(1)
			go func() {
				s.ipv6()
				wg.Done()
			}()
		case "kernel":
			newOrder = append(newOrder, "Kernel")

			wg.Add(1)
			go func() {
				s.kernel()
				wg.Done()
			}()
		case "os":
			newOrder = append(newOrder, "OS")

			wg.Add(1)
			go func() {
				s.operatingSystem()
				wg.Done()
			}()
		case "ram":
			newOrder = append(newOrder, "RAM")

			wg.Add(1)
			go func() {
				s.ram()
				wg.Done()
			}()
		case "shell":
			newOrder = append(newOrder, "Shell")

			wg.Add(1)
			go func() {
				s.shell()
				wg.Done()
			}()
		case "tty":
			newOrder = append(newOrder, "TTY")

			wg.Add(1)
			go func() {
				s.tty()
				wg.Done()
			}()
		case "uptime":
			newOrder = append(newOrder, "Uptime")

			wg.Add(1)
			go func() {
				s.uptime()
				wg.Done()
			}()
		}
	}

	s.order = newOrder
	wg.Wait()
	s.calcSize()
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
		s.CPU = ""
		return s.CPU
	}

	r = regexp.MustCompile(`(cpu model|model name)\s+:\s+(.+)`)
	for _, match := range r.FindAllStringSubmatch(string(info), -1) {
		s.CPU = match[2]
		break
	}

	r = regexp.MustCompile(`\((R|TM)\)| (@|CPU)`)
	s.CPU = r.ReplaceAllString(s.CPU, "")

	r = regexp.MustCompile(`\s+`)
	s.CPU = r.ReplaceAllString(s.CPU, " ")

	return s.CPU
}

func (s *SysInfo) exec(cmd string, cli ...string) string {
	var e error
	var o []byte

	if (cmd == "") || (where.Is(cmd) == "") {
		return ""
	}

	if o, e = exec.Command(cmd, cli...).Output(); e != nil {
		return ""
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

func (s *SysInfo) format(k string, v string, max int) string {
	var filler string
	var line string

	for i := 0; i < max-len(k); i++ {
		filler += " "
	}

	line = " " + hl.Hilights(s.fieldColors, filler+k+":")
	line += " "
	line += hl.Hilights(s.dataColors, v)

	return line
}

func (s *SysInfo) fsUsage(path string) string {
	var usage string
	var words []string

	usage = s.exec("df", "-h", path)

	for _, line := range strings.Split(usage, "\n") {
		words = strings.Fields(line)
		if (len(words) == 6) && (words[5] == path) {
			return words[2] + " / " + words[1] + " (" + words[4] + ")"
		}
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
			s.OS = ""
			return s.OS
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

// SetDataColors will set the color values for the field data. See
// gitlab.com/mjwhitta/hilighter for valid colors.
func (s *SysInfo) SetDataColors(colors ...string) {
	s.dataColors = colors
}

// SetFieldColors will set the color values for the field names. See
// gitlab.com/mjwhitta/hilighter for valid colors.
func (s *SysInfo) SetFieldColors(colors ...string) {
	s.fieldColors = colors
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
	var hasKey bool
	var max int
	var out []string
	var tmp []byte

	tmp, _ = json.Marshal(s)
	json.Unmarshal(tmp, &data)

	for k := range data {
		if len(k) > max {
			max = len(k)
		}
	}

	for _, field := range s.order {
		switch field {
		case "Blank":
			out = append(out, "")
		case "Colors":
			if len(s.Colors) > 0 {
				out = append(out, " "+s.Colors)
			}
		case "FS":
			field = "RootFS"
			if _, hasKey = data[field]; hasKey {
				out = append(out, s.format(field, data[field], max))
			}

			field = "HomeFS"
			if _, hasKey = data[field]; hasKey {
				out = append(out, s.format(field, data[field], max))
			}
		default:
			if _, hasKey = data[field]; hasKey {
				out = append(out, s.format(field, data[field], max))
			}
		}
	}

	return strings.Join(out, "\n")
}

func (s *SysInfo) tty() string {
	var e error
	if s.TTY, e = os.Readlink("/proc/self/fd/0"); e != nil {
		s.TTY = ""
	}
	return s.TTY
}

func (s *SysInfo) uptime() string {
	var r *regexp.Regexp

	s.Uptime = s.exec("uptime")

	// Fail fast
	if s.Uptime == "" {
		return s.Uptime
	}

	// Strip extra whitespace
	r = regexp.MustCompile(`\s+`)
	s.Uptime = r.ReplaceAllString(s.Uptime, " ")

	// Strip leading and trailing data
	r = regexp.MustCompile(`^.*up\s+|,\s+\d+\s+user.+$`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "")

	// Convert hours:mins to match days
	r = regexp.MustCompile(`0?(\d+):0?(\d+)`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "$1 hour, $2 min")

	// Remove 0 hours and mins
	r = regexp.MustCompile(`(^|,\s+)0\s+(hour|min)`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "")

	// Make plural if needed
	r = regexp.MustCompile(`((\d\d+|[2-9])\s+(hour|min))`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "${1}s")

	// Remove leading and trailing commas or whitespace
	r = regexp.MustCompile(`^(,?\s*)+|(,?\s*)+$`)
	s.Uptime = r.ReplaceAllString(s.Uptime, "")

	return s.Uptime
}

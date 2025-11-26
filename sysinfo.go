package sysinfo

import (
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"

	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/where"
)

// SysInfo is a struct containing relevant system information.
type SysInfo struct {
	Colors string   `json:"-"`
	CPU    string   `json:"cpu,omitempty"`
	Height int      `json:"-"`
	HomeFS string   `json:"homefs,omitempty"`
	Host   string   `json:"host,omitempty"`
	IPv4   []string `json:"ipv4,omitempty"`
	IPv6   []string `json:"ipv6,omitempty"`
	Kernel string   `json:"kernel,omitempty"`
	OS     string   `json:"os,omitempty"`
	RAM    string   `json:"ram,omitempty"`
	RootFS string   `json:"rootfs,omitempty"`
	Shell  string   `json:"shell,omitempty"`
	TTY    string   `json:"tty,omitempty"`
	Uptime string   `json:"uptime,omitempty"`
	Width  int      `json:"-"`

	dataColors  []string
	fieldColors []string
	ipMutex     *sync.Mutex
	ips         map[string][]string
	order       []string
}

// New will return a SysInfo pointer. A list of fields can be
// supplied if all info is not wanted.
func New(fields ...string) *SysInfo {
	var s *SysInfo = &SysInfo{ipMutex: &sync.Mutex{}}

	s.order = fields
	if len(fields) == 0 {
		s.order = []string{
			"host",
			"os",
			"kernel",
			"uptime",
			"ip",
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
	s.ips = nil
	s.IPv4 = []string{}
	s.IPv6 = []string{}
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
	var collectFuncs map[string]func() = map[string]func(){
		"blank":  nil,
		"colors": s.colors,
		"cpu":    s.cpu,
		"fs":     s.filesystems,
		"host":   s.hostname,
		"ip":     s.ipAddresses,
		"kernel": s.kernel,
		"os":     s.operatingSystem,
		"ram":    s.ram,
		"shell":  s.shell,
		"tty":    s.tty,
		"uptime": s.uptime,
	}
	var newOrder []string
	var wg sync.WaitGroup

	for _, field := range s.order {
		field = strings.ToLower(field)

		if collect, ok := collectFuncs[field]; ok {
			newOrder = append(newOrder, field)

			if collect == nil {
				continue
			}

			wg.Add(1)

			go func(f func()) {
				f()
				wg.Done()
			}(collect)
		}
	}

	s.order = newOrder

	wg.Wait()
	s.calcSize()
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

func (s *SysInfo) format(k string, v string, maxWidth int) string {
	var filler string = strings.Repeat(" ", maxWidth-len(k)+1)
	var sb strings.Builder

	sb.WriteString(filler)
	sb.WriteString(hl.Hilights(s.fieldColors, k+":"))
	sb.WriteString(" ")
	sb.WriteString(hl.Hilights(s.dataColors, v))

	return sb.String()
}

func (s *SysInfo) getIPs() map[string][]string {
	var addrs []net.Addr
	var e error
	var ifaces []net.Interface
	var ip net.IP
	var tmp string

	s.ipMutex.Lock()
	defer s.ipMutex.Unlock()

	if (s.ips != nil) && (len(s.ips) > 0) {
		return s.ips
	}

	s.ips = map[string][]string{}

	if ifaces, e = net.Interfaces(); e != nil {
		return s.ips
	}

	for _, iface := range ifaces {
		if addrs, e = iface.Addrs(); e != nil {
			continue
		}

		for _, addr := range addrs {
			tmp = addr.String()
			tmp = tmp[0:strings.Index(tmp, "/")]

			if ip = net.ParseIP(tmp); (ip == nil) || ip.IsLoopback() {
				continue
			}

			if ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
				continue
			}

			if strings.HasPrefix(iface.Name, "docker") {
				continue
			}

			s.ips[iface.Name] = append(
				s.ips[iface.Name],
				addr.String(),
			)
		}
	}

	return s.ips
}

func (s *SysInfo) hostname() {
	s.Host = ""
	if host, e := os.Hostname(); e == nil {
		s.Host = strings.TrimSpace(host)
	}
}

func (s *SysInfo) ipAddresses() {
	s.ipv4()
	s.ipv6()
}

func (s *SysInfo) ipv4() {
	var ip net.IP
	var tmp string

	for iface, ips := range s.getIPs() {
		for _, v := range ips {
			tmp = v[0:strings.Index(v, "/")]

			if ip = net.ParseIP(tmp); ip == nil {
				continue
			}

			if ip.To4() != nil {
				s.IPv4 = append(s.IPv4, iface+" "+v)
			}
		}
	}

	sort.Strings(s.IPv4)
}

func (s *SysInfo) ipv6() {
	var ip net.IP
	var tmp string

	for iface, ips := range s.getIPs() {
		for _, v := range ips {
			tmp = v[0:strings.Index(v, "/")]

			if ip = net.ParseIP(tmp); ip == nil {
				continue
			}

			if ip.To4() == nil {
				s.IPv6 = append(s.IPv6, iface+" "+v)
			}
		}
	}

	sort.Strings(s.IPv6)
}

// SetDataColors will set the color values for the field data. See
// github.com/mjwhitta/hilighter for valid colors.
func (s *SysInfo) SetDataColors(colors ...string) {
	s.dataColors = colors
}

// SetFieldColors will set the color values for the field names. See
// github.com/mjwhitta/hilighter for valid colors.
func (s *SysInfo) SetFieldColors(colors ...string) {
	s.fieldColors = colors
}

// String will return a string representation of the SysInfo.
func (s *SysInfo) String() string {
	var data map[string]string = map[string]string{}
	var maxWidth int
	var out []string
	var tmp []byte

	tmp, _ = json.Marshal(s)
	_ = json.Unmarshal(tmp, &data)

	for k := range data {
		if len(k) > maxWidth {
			maxWidth = len(k)
		}
	}

	for _, field := range s.order {
		switch field {
		case "blank":
			out = append(out, "")
		case "colors":
			if s.Colors != "" {
				out = append(out, " "+s.Colors)
			}
		case "fs":
			field = "rootfs"
			if _, ok := data[field]; ok {
				out = append(
					out,
					s.format(titleCase[field], data[field], maxWidth),
				)
			}

			field = "homefs"
			if _, ok := data[field]; ok {
				out = append(
					out,
					s.format(titleCase[field], data[field], maxWidth),
				)
			}
		case "ip":
			field = "ipv4"

			for _, ip := range s.IPv4 {
				out = append(
					out,
					s.format(titleCase[field], ip, maxWidth),
				)
			}

			field = "ipv6"

			for _, ip := range s.IPv6 {
				out = append(
					out,
					s.format(titleCase[field], ip, maxWidth),
				)
			}
		default:
			if _, ok := data[field]; ok {
				out = append(
					out,
					s.format(titleCase[field], data[field], maxWidth),
				)
			}
		}
	}

	return strings.Join(out, "\n")
}

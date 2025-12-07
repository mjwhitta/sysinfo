//go:build windows

package sysinfo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/windows/registry"
)

type clientID struct {
	UniqueProcess uintptr
	UniqueThread  uintptr
}

type objectAttrs struct {
	Length                   uintptr
	RootDirectory            uintptr
	ObjectName               uintptr
	Attributes               uintptr
	SecurityDescriptor       uintptr
	SecurityQualityOfService uintptr
}

func (s *SysInfo) colors() {
	// Needs hilighter support
	s.Colors = ""
}

func (s *SysInfo) cpu() {
	var cpu string
	var e error
	var k registry.Key

	s.CPU = "unknown"

	k, e = registry.OpenKey(
		registry.LOCAL_MACHINE,
		filepath.Join(
			"Hardware",
			"Description",
			"System",
			"CentralProcessor",
			"0",
		),
		registry.QUERY_VALUE,
	)
	if e != nil {
		return
	}
	defer func() {
		_ = k.Close()
	}()

	if cpu, _, e = k.GetStringValue("ProcessorNameString"); e != nil {
		return
	}

	s.CPU = reCPUBrand.ReplaceAllString(cpu, "")
	s.CPU = reWhiteSpace.ReplaceAllString(s.CPU, " ")
}

func (s *SysInfo) filesystems() {
	var home string = strings.ToLower(os.Getenv("HOMEDRIVE"))

	if s.RootFS = s.fsUsage("c:"); s.RootFS == "" {
		s.RootFS = "unknown"
	}

	if home != "c:" {
		s.HomeFS = s.fsUsage(home)
	}
}

func (s *SysInfo) fsUsage(path string) string {
	var cmds []string = []string{
		fmt.Sprintf(
			"gcim win32_logicaldisk -filter \"name='%s'\"",
			path,
		),
		"select deviceid,freespace,size",
	}
	var cols []string
	var e error
	var free int
	var mb int = 1024 * 1024 * 1024
	var total int
	var usage string = s.exec(
		"powershell",
		"-c",
		strings.Join(cmds, "|"),
	)
	var used int

	path = strings.ToLower(path)

	for _, line := range strings.Split(usage, "\n") {
		cols = strings.Fields(strings.ToLower(line))

		//nolint:mnd // Validate output format
		if (len(cols) == 3) && (cols[0] == path) {
			if free, e = strconv.Atoi(cols[1]); e != nil {
				return s.RAM
			}

			if total, e = strconv.Atoi(cols[2]); e != nil {
				return s.RAM
			}

			free /= mb
			total /= mb
			used = total - free

			return fmt.Sprintf(
				"%dG / %dG (%d%%)",
				used,
				total,
				100*used/total,
			)
		}
	}

	return ""
}

func (s *SysInfo) kernel() {
	var build string
	var e error
	var k registry.Key
	var kernel string
	var minor uint64

	s.Kernel = "unknown"

	k, e = registry.OpenKey(
		registry.LOCAL_MACHINE,
		filepath.Join(
			"Software",
			"Microsoft",
			"Windows NT",
			"CurrentVersion",
		),
		registry.QUERY_VALUE,
	)
	if e != nil {
		return
	}
	defer func() {
		_ = k.Close()
	}()

	if kernel, _, e = k.GetStringValue("DisplayVersion"); e != nil {
		return
	}

	if build, _, e = k.GetStringValue("CurrentBuild"); e != nil {
		return
	}

	s.Kernel = kernel + " (OS Build " + build

	if minor, _, e = k.GetIntegerValue("UBR"); e == nil {
		s.Kernel += fmt.Sprintf(".%d", minor)
	}

	s.Kernel += ")"
}

func (s *SysInfo) operatingSystem() {
	var e error
	var k registry.Key
	var os string

	s.OS = "Windows"

	k, e = registry.OpenKey(
		registry.LOCAL_MACHINE,
		filepath.Join(
			"Software",
			"Microsoft",
			"Windows NT",
			"CurrentVersion",
		),
		registry.QUERY_VALUE,
	)
	if e != nil {
		return
	}
	defer func() {
		_ = k.Close()
	}()

	if os, _, e = k.GetStringValue("ProductName"); e != nil {
		return
	}

	s.OS = os
}

func (s *SysInfo) ram() {
	var cmds []string
	var e error
	var free int
	var mb int = 1024 * 1024
	var out string
	var total int

	s.RAM = "unknown"

	cmds = []string{
		"get-counter \"\\memory\\available bytes\"",
		"select -expand countersamples",
		"select -expand cookedvalue",
	}
	out = s.exec(
		"powershell",
		"-c",
		strings.Join(cmds, "|"),
	)

	if free, e = strconv.Atoi(out); e != nil {
		return
	}

	cmds = []string{
		"gcim win32_physicalmemory",
		"measure -property capacity -sum",
		"select -expand sum",
	}
	out = s.exec(
		"powershell",
		"-c",
		strings.Join(cmds, "|"),
	)

	if total, e = strconv.Atoi(out); e != nil {
		return
	}

	s.RAM = fmt.Sprintf("%d MB / %d MB", (total-free)/mb, total/mb)
}

func (s *SysInfo) shell() {
	var sh string

	s.Shell = "unknown"

	sh = s.exec(
		"powershell",
		"-c",
		fmt.Sprintf("(get-process -id %d).processname", os.Getppid()),
	)
	if sh != "" {
		s.Shell = sh
	}
}

func (s *SysInfo) tty() {
	s.TTY = ""
}

func (s *SysInfo) uptime() {
	var out string = s.exec(
		"powershell",
		"-c",
		"(date) - (gcim win32_operatingsystem).lastbootuptime",
	)
	var stop bool
	var unit string

	for _, line := range strings.Split(out, "\n") {
		unit = ""

		switch {
		case strings.HasPrefix(line, "Days"):
			unit = "day"
		case strings.HasPrefix(line, "Hours"):
			unit = "hour"
		case strings.HasPrefix(line, "Minutes"):
			unit = "min"
			stop = true
		}

		if unit == "" {
			continue
		}

		//nolint:mnd // Unit : value == 3 fields
		if tmp := strings.Fields(line); len(tmp) == 3 {
			if tmp[2] != "0" {
				if s.Uptime != "" {
					s.Uptime += ", "
				}

				s.Uptime += tmp[2] + " " + unit
				if tmp[2] != "1" {
					s.Uptime += "s"
				}
			}
		}

		if stop {
			break
		}
	}

	if s.Uptime == "" {
		s.Uptime = "0 mins"
	}
}

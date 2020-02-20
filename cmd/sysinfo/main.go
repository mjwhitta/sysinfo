package main

import (
	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/sysinfo"
)

// Exit status
const (
	Good            int = 0
	InvalidOption   int = 1
	InvalidArgument int = 2
	ExtraArguments  int = 3
	Exception       int = 4
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r.(error).Error())
			}
			errx(Exception, r.(error).Error())
		}
	}()

	var dataColors []string
	var fieldColors []string
	var s *sysinfo.SysInfo

	validate()

	s = sysinfo.New()

	dataColors, _ = config.GetStringArray("dataColors")
	s.SetDataColors(dataColors...)

	fieldColors, _ = config.GetStringArray("fieldColors")
	s.SetFieldColors(fieldColors...)

	hl.Println(s)
}

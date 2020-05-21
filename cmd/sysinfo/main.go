package main

import (
	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/log"
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
			log.ErrX(Exception, r.(error).Error())
		}
	}()

	var s *sysinfo.SysInfo

	validate()

	s = sysinfo.New(flags.fields...)
	s.SetDataColors(config.GetStringArray("dataColors")...)
	s.SetFieldColors(config.GetStringArray("fieldColors")...)

	if s.String() != "" {
		hl.Println(s)
	}
}

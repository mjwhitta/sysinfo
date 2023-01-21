package main

import (
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/log"
	"github.com/mjwhitta/sysinfo"
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

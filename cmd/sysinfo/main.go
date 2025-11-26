package main

import (
	"fmt"

	"github.com/mjwhitta/log"
	"github.com/mjwhitta/sysinfo"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r)
			}

			switch r := r.(type) {
			case error:
				log.ErrX(Exception, r.Error())
			case string:
				log.ErrX(Exception, r)
			}
		}
	}()

	var s *sysinfo.SysInfo

	validate()

	s = sysinfo.New(flags.fields...)
	s.SetDataColors(cfg.DataColors...)
	s.SetFieldColors(cfg.FieldColors...)

	if s.String() != "" {
		fmt.Println(s)
	}
}

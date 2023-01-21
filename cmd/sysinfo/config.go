package main

import "github.com/mjwhitta/jsoncfg"

var config *jsoncfg.JSONCfg

func init() {
	config = jsoncfg.New("~/.config/sysinfo/rc")
	config.SetDefault([]string{"green"}, "dataColors")
	config.SetDefault([]string{"blue"}, "fieldColors")
	config.SaveDefault()
	config.Reset()
}

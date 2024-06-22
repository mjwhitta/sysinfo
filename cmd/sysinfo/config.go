package main

import "github.com/mjwhitta/jsoncfg"

var config *jsoncfg.JSONCfg

func init() {
	config = jsoncfg.New("~/.config/sysinfo/rc")
	_ = config.SetDefault([]string{"green"}, "dataColors")
	_ = config.SetDefault([]string{"blue"}, "fieldColors")
	_ = config.SaveDefault()
	_ = config.Reset()
}

package main

import "gitlab.com/mjwhitta/jsoncfg"

var config *jsoncfg.JSONCfg

func init() {
	config = jsoncfg.New("~/.config/sysinfo/rc")
	config.SetDefault("dataColors", []string{"green", "on_default"})
	config.SetDefault("fieldColors", []string{"blue", "on_default"})
	config.SaveDefault()
	config.Reset()
}

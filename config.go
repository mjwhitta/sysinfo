package sysinfo

import "gitlab.com/mjwhitta/jsoncfg"

var config *jsoncfg.JSONCfg

func init() {
	config = jsoncfg.New("~/.config/sysinfo/rc")
	config.SetDefault("kbg", "on_default")
	config.SetDefault("kfg", "blue")
	config.SetDefault("vbg", "on_default")
	config.SetDefault("vfg", "green")
	config.SaveDefault()
	config.Reset()
}

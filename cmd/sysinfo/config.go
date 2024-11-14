package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type config struct {
	DataColors  []string `json:"dataColors"`
	FieldColors []string `json:"fieldColors"`
}

var cfg *config

func init() {
	var b []byte
	var e error
	var fn string

	if fn, e = os.UserConfigDir(); e != nil {
		panic(fmt.Errorf("user has no cfg directory: %w", e))
	}

	fn = filepath.Join(fn, "sysinfo", "rc")
	b, e = os.ReadFile(fn)

	if (e != nil) || (len(bytes.TrimSpace(b)) == 0) {
		// Default cfg
		cfg = &config{
			DataColors:  []string{"green"},
			FieldColors: []string{"blue"},
		}

		b, _ = json.MarshalIndent(&cfg, "", "  ")

		if e = os.MkdirAll(filepath.Dir(fn), 0o700); e != nil {
			e = fmt.Errorf(
				"failed to create directory %s: %w",
				filepath.Dir(fn),
				e,
			)
			panic(e)
		}

		if e = os.WriteFile(fn, append(b, '\n'), 0o600); e != nil {
			panic(fmt.Errorf("failed to write %s: %w", fn, e))
		}
	} else {
		if e = json.Unmarshal(b, &cfg); e != nil {
			panic(fmt.Errorf("invalid cfg: %w", e))
		}
	}

	if cfg.DataColors == nil {
		cfg.DataColors = []string{"green"}
	}

	if cfg.FieldColors == nil {
		cfg.FieldColors = []string{"blue"}
	}
}

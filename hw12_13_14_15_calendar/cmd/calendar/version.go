package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

var version = struct {
	Release   string
	BuildDate string
	GitHash   string
}{
	Release:   release,
	BuildDate: buildDate,
	GitHash:   gitHash,
}

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(version); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}

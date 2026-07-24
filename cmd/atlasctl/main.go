package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Iron-Signal-Systems/atlas/internal/config"
)

func main() {
	if len(os.Args) != 3 || os.Args[1] != "validate-config" {
		fmt.Fprintln(os.Stderr, "usage: atlasctl validate-config <config.json>")
		os.Exit(2)
	}

	data, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config: %v\n", err)
		os.Exit(1)
	}

	var cfg config.File
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "parse config: %v\n", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("configuration is valid")
}

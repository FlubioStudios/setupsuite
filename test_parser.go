package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"suite/suite/config"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run test_parser.go <config_file>")
		os.Exit(1)
	}

	configFile := os.Args[1]
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.ParseConfig(string(content))
	if err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("=== Parsed Configuration ===\n")
	if cfg.SetupSecure != nil {
		fmt.Printf("SSH User: '%s'\n", cfg.SetupSecure.SSHUser)
		fmt.Printf("SSH Key: '%s'\n", cfg.SetupSecure.UserSSHRSA[:min(50, len(cfg.SetupSecure.UserSSHRSA))]+"...")
		fmt.Printf("SSH Port: %d\n", cfg.SetupSecure.SSHPort)

		if cfg.SetupSecure.Config != nil {
			fmt.Printf("Config Type: '%s'\n", cfg.SetupSecure.Config.Type)
			fmt.Printf("Domain: '%s'\n", cfg.SetupSecure.Config.Domain)
			fmt.Printf("Email: '%s'\n", cfg.SetupSecure.Config.Email)
		}
	}

	if cfg.InstallTools != nil {
		fmt.Printf("Tools: %v\n", cfg.InstallTools.Tools)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

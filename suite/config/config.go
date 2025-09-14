package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadConfig(configPath string) (*ServerConfig, error) {
	path := "/etc/setupsuite"
	config := configPath

	if configPath == "" {
		config = path + "/config.sscfg"
	}

	// Create config directory if it doesn't exist
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		fmt.Println("creating config dir")
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	// Create default config file if it doesn't exist
	if _, err := os.Stat(config); errors.Is(err, os.ErrNotExist) {
		fmt.Println("creating default config")
		err := CreateDefaultConfig("basic", config)
		if err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		fmt.Printf("Default config created at %s. Please edit it and run again.\n", config)
		os.Exit(0)
	}

	fmt.Println("Reading config from:", config)
	dat, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the configuration
	serverConfig, err := ParseConfig(string(dat))
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	fmt.Println("Configuration loaded successfully")
	return serverConfig, nil
}

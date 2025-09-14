package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"suite/suite/config"
)

func main() {
	// Parse command line arguments
	var (
		configPath   = flag.String("config", "/etc/setupsuite/config.sscfg", "Path to configuration file")
		serverType   = flag.String("type", "", "Server type (web, database, docker, proxy, build)")
		generateOnly = flag.Bool("generate", false, "Generate default config and exit")
		verbose      = flag.Bool("verbose", false, "Enable verbose logging of all file operations and command outputs")
		help         = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	// Initialize logging
	if err := InitLogger(*verbose); err != nil {
		fmt.Printf("Warning: Could not initialize logging: %v\n", err)
	}
	defer CloseLogger()

	if *help {
		showHelp()
		return
	}

	if runtime.GOOS == "windows" {
		fmt.Println("The LINUX setup suite is only to be used on linux based systems")
		return
	}

	// Generate default config if requested
	if *generateOnly {
		if *serverType == "" {
			fmt.Println("Please specify server type with -type flag when generating config")
			fmt.Println("Available types: web, database, docker, proxy, build")
			return
		}

		err := config.CreateDefaultConfig(*serverType, *configPath)
		if err != nil {
			fmt.Printf("Error creating config: %v\n", err)
			return
		}

		fmt.Printf("Default %s server config created at %s\n", *serverType, *configPath)
		fmt.Println("Please edit the config file and run the setup again without -generate flag")
		return
	}

	fmt.Println("Starting Serversetup...")
	execute(*configPath)
}

func showHelp() {
	fmt.Println("SetupSuite - Automated Linux Server Setup Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  setupsuite [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -config string    Path to configuration file (default: /etc/setupsuite/config.sscfg)")
	fmt.Println("  -type string      Server type for config generation (web, database, docker, proxy, build)")
	fmt.Println("  -generate         Generate default config file and exit")
	fmt.Println("  -verbose          Enable verbose logging of all file operations and command outputs")
	fmt.Println("  -help             Show this help message")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  setupsuite -generate -type web                    # Generate web server config")
	fmt.Println("  setupsuite -config /path/to/custom.sscfg          # Use custom config")
	fmt.Println("  setupsuite -verbose                               # Run with verbose logging")
	fmt.Println("  setupsuite                                        # Use default config")
	fmt.Println("")
	fmt.Println("Server Types:")
	fmt.Println("  web       - Web server with Nginx and SSL")
	fmt.Println("  database  - Database server (MySQL/PostgreSQL)")
	fmt.Println("  docker    - Docker host with optimized settings")
	fmt.Println("  proxy     - Reverse proxy server")
	fmt.Println("  build     - Build/CI server with dev tools")
}

func execute(configPath string) {
	// Check if running as root
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	if user.Username != "root" {
		fmt.Println("The Setupsuite can only be run as root")
		os.Exit(1)
	}

	// Detect system information
	fmt.Println("Detecting system information...")
	distro, version, err := DetectDistribution()
	if err != nil {
		fmt.Printf("Warning: Could not detect distribution: %v\n", err)
	} else {
		fmt.Printf("Detected: %s %s\n", distro, version)
	}

	// Detect package manager
	pm, err := DetectPackageManager()
	if err != nil {
		fmt.Printf("Warning: Could not detect package manager: %v\n", err)
	} else {
		fmt.Printf("Package manager: %s\n", pm.GetName())
	}

	// Read and parse configuration
	serverConfig, err := config.ReadConfig(configPath)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}

	// Perform server setup based on config
	err = setupServer(serverConfig)
	if err != nil {
		fmt.Printf("Setup failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Server setup completed successfully!")
}

func setupServer(cfg *config.ServerConfig) error {
	// Basic security setup
	if cfg.SetupSecure != nil {
		fmt.Println("Setting up basic security...")
		err := setupBasicSecurity(cfg.SetupSecure)
		if err != nil {
			return fmt.Errorf("security setup failed: %v", err)
		}
	}

	// Install packages
	if cfg.InstallTools != nil && len(cfg.InstallTools.Tools) > 0 {
		err := InstallPackages(cfg.InstallTools.Tools)
		if err != nil {
			return fmt.Errorf("package installation failed: %v", err)
		}
	}

	// Configure firewall
	if cfg.SetupSecure != nil && cfg.SetupSecure.Firewall != nil {
		err := ConfigureFirewall(cfg.SetupSecure.Firewall.OpenPorts)
		if err != nil {
			return fmt.Errorf("firewall configuration failed: %v", err)
		}
	}

	// Server-specific setup
	if cfg.SetupSecure != nil && cfg.SetupSecure.Config != nil {
		serverSetup := &ServerSetup{Config: cfg}

		switch cfg.SetupSecure.Config.Type {
		case config.ServerTypeWeb:
			return serverSetup.SetupWebServer()
		case config.ServerTypeDatabase:
			return serverSetup.SetupDatabaseServer()
		case config.ServerTypeDocker:
			return serverSetup.SetupDockerHost()
		case config.ServerTypeProxy:
			return serverSetup.SetupProxyServer()
		case config.ServerTypeBuild:
			return serverSetup.SetupBuildServer()
		default:
			fmt.Printf("Unknown server type: %s\n", cfg.SetupSecure.Config.Type)
		}
	}

	return nil
}

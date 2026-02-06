package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ParseConfig parses the custom configuration format
func ParseConfig(content string) (*ServerConfig, error) {
	config := &ServerConfig{}
	lines := strings.Split(content, "\n")

	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" || strings.HasPrefix(line, "#") {
			i++
			continue
		}

		if strings.HasPrefix(line, ".setup_secure{") {
			setupSecure, nextIndex := parseSetupSecure(lines, i)
			config.SetupSecure = setupSecure
			i = nextIndex
		} else if strings.HasPrefix(line, ".install_tools{") {
			installTools, nextIndex := parseInstallTools(lines, i)
			config.InstallTools = installTools
			i = nextIndex
		} else {
			i++
		}
	}

	return config, nil
}

func parseSetupSecure(lines []string, startIndex int) (*SetupSecure, int) {
	setupSecure := &SetupSecure{}
	i := startIndex + 1

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "}" {
			break
		}

		if strings.Contains(line, ":") && !strings.HasPrefix(line, ".") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes and trailing commas more robustly
			value = strings.Trim(value, " \t\r\n") // Remove whitespace
			if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `",`) {
				value = value[1 : len(value)-2] // Remove quote and comma
			} else if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
				value = value[1 : len(value)-1] // Remove quotes
			}
			value = strings.TrimSuffix(value, ",") // Remove any remaining trailing comma

			switch key {
			case "ssh_user":
				setupSecure.SSHUser = value
			case "user_ssh_rsa":
				setupSecure.UserSSHRSA = value
			case "ssh_port":
				if port, err := strconv.Atoi(value); err == nil {
					setupSecure.SSHPort = port
				}
			}
		} else if strings.HasPrefix(line, ".configuration{") {
			config, nextIndex := parseConfiguration(lines, i)
			setupSecure.Config = config
			i = nextIndex
			continue
		} else if strings.HasPrefix(line, ".firewall{") {
			firewall, nextIndex := parseFirewall(lines, i)
			setupSecure.Firewall = firewall
			i = nextIndex
			continue
		}
		i++
	}

	return setupSecure, i + 1
}

func parseConfiguration(lines []string, startIndex int) (*Config, int) {
	config := &Config{
		Options: make(map[string]string),
	}
	i := startIndex + 1

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "}" || line == "}," {
			break
		}

		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes and trailing commas more robustly
			value = strings.Trim(value, " \t\r\n") // Remove whitespace
			if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `",`) {
				value = value[1 : len(value)-2] // Remove quote and comma
			} else if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
				value = value[1 : len(value)-1] // Remove quotes
			}
			value = strings.TrimSuffix(value, ",") // Remove any remaining trailing comma

			switch key {
			case "type":
				config.Type = value
			case "domain":
				config.Domain = value
			case "email":
				config.Email = value
			default:
				config.Options[key] = value
			}
		}
		i++
	}

	return config, i + 1
}

func parseFirewall(lines []string, startIndex int) (*Firewall, int) {
	firewall := &Firewall{}
	i := startIndex + 1

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "}" {
			break
		}

		if strings.HasPrefix(line, "open_ports:") {
			i++
			// Parse array
			for i < len(lines) {
				arrayLine := strings.TrimSpace(lines[i])
				if arrayLine == "]" {
					break
				}
				if arrayLine != "[" && arrayLine != "" {
					portStr := strings.Trim(arrayLine, ",")
					if port, err := strconv.Atoi(portStr); err == nil {
						firewall.OpenPorts = append(firewall.OpenPorts, port)
					}
				}
				i++
			}
		}
		i++
	}

	return firewall, i + 1
}

func parseInstallTools(lines []string, startIndex int) (*InstallTools, int) {
	installTools := &InstallTools{}
	i := startIndex + 1

	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "}" {
			break
		}

		if strings.HasPrefix(line, "tools:") {
			i++
			// Parse array
			for i < len(lines) {
				arrayLine := strings.TrimSpace(lines[i])
				if arrayLine == "]" {
					break
				}
				if arrayLine != "[" && arrayLine != "" {
					tool := strings.TrimSpace(arrayLine)
					// Remove quotes and trailing commas more robustly
					tool = strings.Trim(tool, " \t\r\n") // Remove whitespace
					if strings.HasPrefix(tool, `"`) && strings.HasSuffix(tool, `",`) {
						tool = tool[1 : len(tool)-2] // Remove quote and comma
					} else if strings.HasPrefix(tool, `"`) && strings.HasSuffix(tool, `"`) {
						tool = tool[1 : len(tool)-1] // Remove quotes
					}
					tool = strings.TrimSuffix(tool, ",") // Remove any remaining trailing comma
					if tool != "" {
						installTools.Tools = append(installTools.Tools, tool)
					}
				}
				i++
			}
		}
		i++
	}

	return installTools, i + 1
}

// CreateDefaultConfig creates default configuration files for different server types
func CreateDefaultConfig(serverType string, configPath string) error {
	var configContent string

	switch serverType {
	case ServerTypeWeb:
		configContent = getWebServerConfig()
	case ServerTypeDatabase:
		configContent = getDatabaseServerConfig()
	case ServerTypeDocker:
		configContent = getDockerHostConfig()
	case ServerTypeProxy:
		configContent = getProxyServerConfig()
	case ServerTypeBuild:
		configContent = getBuildServerConfig()
	default:
		configContent = getBasicServerConfig()
	}

	// Create parent directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		fmt.Printf("Creating config directory: %s\n", configDir)
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("permission denied creating config directory %s (try running with sudo or use a different path with -config flag): %v", configDir, err)
			}
			return fmt.Errorf("failed to create config directory %s: %v", configDir, err)
		}
		fmt.Printf("Config directory created successfully\n")
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(configContent)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func getWebServerConfig() string {
	return `.setup_secure{
	ssh_user: "admin",
	user_ssh_rsa: "REPLACE_WITH_YOUR_SSH_KEY",
	ssh_port: 22022,
	.configuration{
		type: "web",
		domain: "example.com",
		email: "admin@example.com"
	},
	.firewall{
		open_ports: [
			22022,
			80,
			443
		]
	}
}

.install_tools{
	tools: [
		"nginx",
		"certbot",
		"python3-certbot-nginx",
		"ufw",
		"htop",
		"curl",
		"git"
	]
}`
}

func getDatabaseServerConfig() string {
	return `.setup_secure{
	ssh_user: "dbadmin",
	user_ssh_rsa: "REPLACE_WITH_YOUR_SSH_KEY",
	ssh_port: 22022,
	.configuration{
		type: "database",
		db_engine: "mysql",
		root_password: "REPLACE_WITH_SECURE_PASSWORD"
	},
	.firewall{
		open_ports: [
			22022,
			3306
		]
	}
}

.install_tools{
	tools: [
		"mysql-server",
		"mysql-client",
		"ufw",
		"htop",
		"curl"
	]
}`
}

func getDockerHostConfig() string {
	return `.setup_secure{
	ssh_user: "dockeradmin",
	user_ssh_rsa: "REPLACE_WITH_YOUR_SSH_KEY",
	ssh_port: 22022,
	.configuration{
		type: "docker"
	},
	.firewall{
		open_ports: [
			22022,
			80,
			443,
			2376
		]
	}
}

.install_tools{
	tools: [
		"docker.io",
		"docker-compose",
		"ufw",
		"htop",
		"curl",
		"git"
	]
}`
}

func getProxyServerConfig() string {
	return `.setup_secure{
	ssh_user: "proxyadmin",
	user_ssh_rsa: "REPLACE_WITH_YOUR_SSH_KEY",
	ssh_port: 22022,
	.configuration{
		type: "proxy",
		domain: "proxy.example.com",
		email: "admin@example.com"
	},
	.firewall{
		open_ports: [
			22022,
			80,
			443
		]
	}
}

.install_tools{
	tools: [
		"nginx",
		"certbot",
		"python3-certbot-nginx",
		"ufw",
		"htop",
		"curl"
	]
}`
}

func getBuildServerConfig() string {
	return `.setup_secure{
	ssh_user: "buildadmin",
	user_ssh_rsa: "REPLACE_WITH_YOUR_SSH_KEY",
	ssh_port: 22022,
	.configuration{
		type: "build"
	},
	.firewall{
		open_ports: [
			22022,
			80,
			443,
			8080
		]
	}
}

.install_tools{
	tools: [
		"git",
		"nodejs",
		"npm",
		"python3",
		"python3-pip",
		"build-essential",
		"docker.io",
		"ufw",
		"htop",
		"curl"
	]
}`
}

func getBasicServerConfig() string {
	return `.setup_secure{
	ssh_user: "admin",
	user_ssh_rsa: "REPLACE_WITH_YOUR_SSH_KEY",
	ssh_port: 22022,
	.configuration{
		type: "basic"
	},
	.firewall{
		open_ports: [
			22022
		]
	}
}

.install_tools{
	tools: [
		"ufw",
		"htop",
		"curl",
		"git",
		"nano"
	]
}`
}

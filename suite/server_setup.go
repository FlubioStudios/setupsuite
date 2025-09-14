package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"suite/suite/config"
)

// ServerSetup handles different server setup types
type ServerSetup struct {
	Config *config.ServerConfig
}

// SetupWebServer configures a web server with Nginx and SSL
func (s *ServerSetup) SetupWebServer() error {
	fmt.Println("Setting up Web Server...")

	// Get service manager
	sm, err := NewServiceManager()
	if err != nil {
		fmt.Printf("Warning: Could not detect service manager: %v\n", err)
		return err
	}

	// Install and configure Nginx
	if err := sm.Enable("nginx"); err != nil {
		fmt.Printf("Warning: Could not enable nginx: %v\n", err)
	}
	if err := sm.Start("nginx"); err != nil {
		fmt.Printf("Warning: Could not start nginx: %v\n", err)
	}

	// Configure basic nginx site
	if s.Config.SetupSecure.Config.Domain != "" {
		s.setupNginxSite(s.Config.SetupSecure.Config.Domain)
	}

	// Setup SSL if domain and email are provided
	if s.Config.SetupSecure.Config.Domain != "" && s.Config.SetupSecure.Config.Email != "" {
		s.setupSSL(s.Config.SetupSecure.Config.Domain, s.Config.SetupSecure.Config.Email)
	}

	return nil
}

// SetupDatabaseServer configures a database server
func (s *ServerSetup) SetupDatabaseServer() error {
	fmt.Println("Setting up Database Server...")

	dbEngine := s.Config.SetupSecure.Config.Options["db_engine"]
	if dbEngine == "" {
		dbEngine = "mysql"
	}

	switch dbEngine {
	case "mysql":
		return s.setupMySQL()
	case "postgresql":
		return s.setupPostgreSQL()
	default:
		return s.setupMySQL()
	}
}

// SetupDockerHost configures a Docker host
func (s *ServerSetup) SetupDockerHost() error {
	fmt.Println("Setting up Docker Host...")

	// Get service manager
	sm, err := NewServiceManager()
	if err != nil {
		fmt.Printf("Warning: Could not detect service manager: %v\n", err)
		return err
	}

	// Add user to docker group
	if s.Config.SetupSecure.SSHUser != "" {
		exec.Command("usermod", "-aG", "docker", s.Config.SetupSecure.SSHUser).Run()
	}

	// Configure Docker daemon
	s.setupDockerDaemon()

	// Enable and start Docker
	sm.Enable("docker")
	sm.Start("docker")

	return nil
}

// SetupProxyServer configures a reverse proxy server
func (s *ServerSetup) SetupProxyServer() error {
	fmt.Println("Setting up Proxy Server...")

	// Get service manager
	sm, err := NewServiceManager()
	if err != nil {
		fmt.Printf("Warning: Could not detect service manager: %v\n", err)
		return err
	}

	// Configure Nginx as reverse proxy
	s.setupNginxProxy()

	// Enable and start Nginx
	sm.Enable("nginx")
	sm.Start("nginx")

	// Setup SSL if domain and email are provided
	if s.Config.SetupSecure.Config.Domain != "" && s.Config.SetupSecure.Config.Email != "" {
		s.setupSSL(s.Config.SetupSecure.Config.Domain, s.Config.SetupSecure.Config.Email)
	}

	return nil
}

// SetupBuildServer configures a build/CI server
func (s *ServerSetup) SetupBuildServer() error {
	fmt.Println("Setting up Build Server...")

	// Get service manager
	sm, err := NewServiceManager()
	if err != nil {
		fmt.Printf("Warning: Could not detect service manager: %v\n", err)
	}

	// Add user to docker group for Docker builds
	if s.Config.SetupSecure.SSHUser != "" {
		exec.Command("usermod", "-aG", "docker", s.Config.SetupSecure.SSHUser).Run()
	}

	// Install Node.js LTS
	s.installNodeJS()

	// Setup Python virtual environment tools
	exec.Command("pip3", "install", "virtualenv").Run()

	// Enable Docker if service manager is available
	if sm != nil {
		sm.Enable("docker")
		sm.Start("docker")
	}

	return nil
}

func (s *ServerSetup) setupNginxSite(domain string) {
	nginxConfig := fmt.Sprintf(`server {
    listen 80;
    server_name %s www.%s;
    
    root /var/www/html;
    index index.html index.htm index.nginx-debian.html;
    
    location / {
        try_files $uri $uri/ =404;
    }
}`, domain, domain)

	configPath := "/etc/nginx/sites-available/" + domain
	if file, err := os.Create(configPath); err == nil {
		file.WriteString(nginxConfig)
		file.Close()

		// Enable site
		linkPath := "/etc/nginx/sites-enabled/" + domain
		exec.Command("ln", "-sf", configPath, linkPath).Run()

		// Remove default site
		exec.Command("rm", "-f", "/etc/nginx/sites-enabled/default").Run()

		// Test and reload nginx
		exec.Command("nginx", "-t").Run()

		// Reload nginx using service manager
		if sm, err := NewServiceManager(); err == nil {
			sm.Reload("nginx")
		} else {
			exec.Command("systemctl", "reload", "nginx").Run()
		}
	}
}

func (s *ServerSetup) setupSSL(domain, email string) {
	fmt.Printf("Setting up SSL for %s...\n", domain)

	// Use certbot to get SSL certificate
	cmd := exec.Command("certbot", "--nginx", "-d", domain, "-d", "www."+domain,
		"--non-interactive", "--agree-tos", "--email", email, "--redirect")
	cmd.Run()
}

func (s *ServerSetup) setupMySQL() error {
	fmt.Println("Configuring MySQL...")

	// Get service manager
	sm, err := NewServiceManager()
	if err != nil {
		fmt.Printf("Warning: Could not detect service manager: %v\n", err)
		return err
	}

	// Detect distribution to handle MySQL service names
	distro, _, _ := DetectDistribution()
	serviceName := "mysql"

	switch distro {
	case "rhel", "centos", "fedora":
		serviceName = "mysqld"
	case "alpine":
		serviceName = "mariadb"
	default:
		serviceName = "mysql"
	}

	// Secure MySQL installation (non-interactive)
	exec.Command("mysql_secure_installation").Run()

	// Enable and start MySQL
	sm.Enable(serviceName)
	sm.Start(serviceName)

	return nil
}

func (s *ServerSetup) setupPostgreSQL() error {
	fmt.Println("Configuring PostgreSQL...")

	// Get service manager
	sm, err := NewServiceManager()
	if err != nil {
		fmt.Printf("Warning: Could not detect service manager: %v\n", err)
		return err
	}

	// Detect distribution to handle PostgreSQL service names
	distro, _, _ := DetectDistribution()
	serviceName := "postgresql"

	switch distro {
	case "rhel", "centos", "fedora":
		serviceName = "postgresql"
	case "alpine":
		serviceName = "postgresql"
	default:
		serviceName = "postgresql"
	}

	// Enable and start PostgreSQL
	sm.Enable(serviceName)
	sm.Start(serviceName)

	return nil
}

func (s *ServerSetup) setupDockerDaemon() {
	fmt.Println("Configuring Docker daemon...")

	// Create docker daemon.json for log rotation (from TODO.md)
	daemonConfig := `{
  "log-driver": "local",
  "log-opts": {
    "max-size": "20m",
    "max-file": "5"
  }
}`

	// Create /etc/docker directory if it doesn't exist
	os.MkdirAll("/etc/docker", 0755)

	if file, err := os.Create("/etc/docker/daemon.json"); err == nil {
		file.WriteString(daemonConfig)
		file.Close()

		// Restart Docker to apply changes
		if sm, err := NewServiceManager(); err == nil {
			sm.Restart("docker")
		} else {
			exec.Command("systemctl", "restart", "docker").Run()
		}
	}
}

func (s *ServerSetup) setupNginxProxy() {
	fmt.Println("Configuring Nginx as reverse proxy...")

	nginxConfig := `# Global configuration
upstream backend {
    server 127.0.0.1:3000;
}

server {
    listen 80;
    server_name _;
    
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}`

	configPath := "/etc/nginx/sites-available/proxy"
	if file, err := os.Create(configPath); err == nil {
		file.WriteString(nginxConfig)
		file.Close()

		// Enable site
		linkPath := "/etc/nginx/sites-enabled/proxy"
		exec.Command("ln", "-sf", configPath, linkPath).Run()

		// Remove default site
		exec.Command("rm", "-f", "/etc/nginx/sites-enabled/default").Run()
	}
}

func (s *ServerSetup) installNodeJS() {
	fmt.Println("Installing Node.js LTS...")

	// Detect distribution for appropriate installation method
	distro, _, err := DetectDistribution()
	if err != nil {
		fmt.Printf("Warning: Could not detect distribution: %v\n", err)
		return
	}

	switch distro {
	case "ubuntu", "debian":
		// Install NodeSource repository for Debian/Ubuntu
		exec.Command("curl", "-fsSL", "https://deb.nodesource.com/setup_lts.x", "|", "bash", "-").Run()
		exec.Command("apt-get", "install", "-y", "nodejs").Run()
	case "rhel", "centos", "fedora":
		// Use NodeSource for RHEL-based systems
		exec.Command("curl", "-fsSL", "https://rpm.nodesource.com/setup_lts.x", "|", "bash", "-").Run()
		pm, _ := DetectPackageManager()
		if pm != nil {
			pm.Install([]string{"nodejs"})
		}
	case "arch":
		// Use official Arch repositories
		pm, _ := DetectPackageManager()
		if pm != nil {
			pm.Install([]string{"nodejs", "npm"})
		}
	case "alpine":
		// Use Alpine repositories
		pm, _ := DetectPackageManager()
		if pm != nil {
			pm.Install([]string{"nodejs", "npm"})
		}
	default:
		fmt.Printf("Warning: Node.js installation not configured for distribution: %s\n", distro)
		fmt.Println("Please install Node.js manually")
	}
}

// InstallPackages installs system packages using the detected package manager
func InstallPackages(packages []string) error {
	if len(packages) == 0 {
		fmt.Println("No packages to install")
		return nil
	}

	// Detect the package manager
	pm, err := DetectPackageManager()
	if err != nil {
		return fmt.Errorf("failed to detect package manager: %v", err)
	}

	fmt.Printf("Installing packages using %s: %s\n", pm.GetName(), strings.Join(packages, ", "))

	// Update package list
	if err := pm.Update(); err != nil {
		fmt.Printf("Warning: Failed to update package list: %v\n", err)
		// Continue anyway, as some package managers don't require explicit updates
	}

	// Install packages
	if err := pm.Install(packages); err != nil {
		return fmt.Errorf("failed to install packages with %s: %v", pm.GetName(), err)
	}

	fmt.Println("Package installation completed successfully")
	return nil
}

// ConfigureFirewall sets up firewall rules using the detected firewall manager
func ConfigureFirewall(ports []int) error {
	if len(ports) == 0 {
		fmt.Println("No firewall ports to configure")
		return nil
	}

	// Detect firewall manager
	fwManager, err := GetFirewallManager()
	if err != nil {
		fmt.Printf("Warning: %v. Skipping firewall configuration.\n", err)
		return nil
	}

	fmt.Printf("Configuring firewall using %s...\n", fwManager)

	switch fwManager {
	case "ufw":
		return configureUFW(ports)
	case "firewalld":
		return configureFirewalld(ports)
	case "iptables":
		return configureIptables(ports)
	default:
		return fmt.Errorf("unsupported firewall manager: %s", fwManager)
	}
}

func configureUFW(ports []int) error {
	// Reset UFW to defaults
	exec.Command("ufw", "--force", "reset").Run()

	// Set default policies
	exec.Command("ufw", "default", "deny", "incoming").Run()
	exec.Command("ufw", "default", "allow", "outgoing").Run()

	// Open specified ports
	for _, port := range ports {
		fmt.Printf("Opening port %d (UFW)\n", port)
		exec.Command("ufw", "allow", fmt.Sprintf("%d", port)).Run()
	}

	// Enable UFW
	exec.Command("ufw", "--force", "enable").Run()
	return nil
}

func configureFirewalld(ports []int) error {
	// Start firewalld if not running
	exec.Command("systemctl", "start", "firewalld").Run()
	exec.Command("systemctl", "enable", "firewalld").Run()

	// Open specified ports
	for _, port := range ports {
		fmt.Printf("Opening port %d (firewalld)\n", port)
		exec.Command("firewall-cmd", "--permanent", "--add-port", fmt.Sprintf("%d/tcp", port)).Run()
	}

	// Reload firewall rules
	exec.Command("firewall-cmd", "--reload").Run()
	return nil
}

func configureIptables(ports []int) error {
	// Flush existing rules (be careful!)
	exec.Command("iptables", "-F").Run()

	// Set default policies
	exec.Command("iptables", "-P", "INPUT", "DROP").Run()
	exec.Command("iptables", "-P", "FORWARD", "DROP").Run()
	exec.Command("iptables", "-P", "OUTPUT", "ACCEPT").Run()

	// Allow loopback
	exec.Command("iptables", "-A", "INPUT", "-i", "lo", "-j", "ACCEPT").Run()

	// Allow established connections
	exec.Command("iptables", "-A", "INPUT", "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT").Run()

	// Open specified ports
	for _, port := range ports {
		fmt.Printf("Opening port %d (iptables)\n", port)
		exec.Command("iptables", "-A", "INPUT", "-p", "tcp", "--dport", fmt.Sprintf("%d", port), "-j", "ACCEPT").Run()
	}

	// Save rules (distribution-specific)
	exec.Command("iptables-save").Run() // Basic save, may need distribution-specific handling
	return nil
}

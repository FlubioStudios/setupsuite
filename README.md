# FlubioStudios SetupSuite

A powerful, automated Linux server setup tool designed to quickly configure and secure fresh server installations. SetupSuite uses a custom configuration language to define server roles and automatically installs, configures, and secures your servers according to best practices.

## 🚀 Features

- **Multi-Server Type Support**: Pre-configured templates for web servers, database servers, Docker hosts, proxy servers, and build servers
- **Distribution Agnostic**: Automatically detects and supports multiple Linux distributions (Ubuntu, Debian, RHEL, CentOS, Fedora, Arch, Alpine, openSUSE)
- **Smart Package Management**: Auto-detects package managers (APT, DNF, YUM, Pacman, Zypper, APK)
- **Universal Service Management**: Supports systemd, SysV Init, and OpenRC
- **Flexible Firewall Support**: Works with UFW, firewalld, and iptables
- **Security First**: Automatic SSH hardening, firewall configuration, and user management
- **Custom DSL**: Easy-to-read configuration language for defining server setup
- **SSL Automation**: Automatic SSL certificate generation with Let's Encrypt
- **Docker Integration**: Optimized Docker setup with logging configuration
- **Extensible**: Easy to add new server types and configurations

## 📋 Supported Server Types

| Type | Description | Includes |
|------|-------------|----------|
| **web** | Web server with Nginx | Nginx, SSL/TLS, firewall rules |
| **database** | Database server | MySQL/PostgreSQL, security hardening |
| **docker** | Container host | Docker, docker-compose, optimized logging |
| **proxy** | Reverse proxy | Nginx proxy config, SSL termination |
| **build** | CI/CD build server | Git, Node.js, Docker, build tools |

## 🐧 Supported Linux Distributions

| Distribution | Package Manager | Service Manager | Firewall | Status |
|-------------|----------------|----------------|----------|---------|
| **Ubuntu** 18.04+ | APT | systemd | UFW | ✅ Fully Supported |
| **Debian** 10+ | APT | systemd | UFW | ✅ Fully Supported |
| **RHEL** 8+ | DNF | systemd | firewalld | ✅ Fully Supported |
| **CentOS** 7+ | YUM/DNF | systemd | firewalld | ✅ Fully Supported |
| **Fedora** 30+ | DNF | systemd | firewalld | ✅ Fully Supported |
| **Arch Linux** | Pacman | systemd | UFW/iptables | ✅ Fully Supported |
| **Alpine Linux** | APK | OpenRC | iptables | ✅ Fully Supported |
| **openSUSE** | Zypper | systemd | firewalld | ✅ Fully Supported |

## 🛠️ Installation

### Quick Install (Recommended)
```bash
curl -fsSL https://raw.githubusercontent.com/FlubioStudios/SetupSuite/main/install.sh | bash
```

### Manual Installation
```bash
# Clone the repository
git clone https://github.com/FlubioStudios/SetupSuite.git
cd SetupSuite

# Build the binary
go build -o setupsuite ./suite

# Install system-wide
sudo mv setupsuite /usr/local/bin/
sudo chmod +x /usr/local/bin/setupsuite
```

## 📖 Quick Start

### 1. Generate a Configuration Template
```bash
# Generate a web server configuration
sudo setupsuite -generate -type web -config /etc/setupsuite/web.sscfg

# Other available types: database, docker, proxy, build
```

### 2. Edit the Configuration
```bash
sudo nano /etc/setupsuite/web.sscfg
```

**Important**: Replace the following placeholders:
- `REPLACE_WITH_YOUR_SSH_KEY` - Your SSH public key
- `example.com` - Your actual domain name
- `admin@example.com` - Your email address

### 3. Run the Setup
```bash
sudo setupsuite -config /etc/setupsuite/web.sscfg
```
preview
## 📝 Configuration Language

SetupSuite uses a simple, readable configuration format:

```
.setup_secure{
    ssh_user: "webadmin",
    user_ssh_rsa: "ssh-rsa AAAAB3NzaC1yc2E...",
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
}
```

## 🔧 Command Line Usage

```bash
# Show help
setupsuite -help

# Generate configuration templates
setupsuite -generate -type web -config /path/to/config.sscfg
setupsuite -generate -type database -config /path/to/db-config.sscfg
setupsuite -generate -type docker -config /path/to/docker-config.sscfg

# Run setup with custom config
setupsuite -config /path/to/config.sscfg

# Run with default config (creates one if missing)
setupsuite
```

## 🏗️ What SetupSuite Does

### Security Hardening
- Creates a non-root user with sudo privileges
- Configures SSH key-based authentication
- Disables password authentication
- Changes SSH port (configurable)
- Sets up UFW firewall with specified rules
- Installs and configures fail2ban

### System Updates
- Updates package repositories
- Upgrades all system packages
- Enables automatic security updates

### Server-Specific Setup

#### Web Server
- Installs and configures Nginx
- Sets up SSL certificates with Let's Encrypt
- Configures basic virtual host
- Optimizes Nginx for performance

#### Database Server
- Installs MySQL or PostgreSQL
- Runs security hardening scripts
- Configures firewall for database ports

#### Docker Host
- Installs Docker and docker-compose
- Configures Docker daemon with log rotation
- Adds user to docker group
- Optimizes Docker settings

#### Proxy Server
- Configures Nginx as reverse proxy
- Sets up SSL termination
- Configures upstream servers

#### Build Server
- Installs development tools
- Sets up Node.js, Python, Java
- Configures Docker for CI/CD
- Installs build dependencies

## 📁 Example Configurations

See the `examples/` directory for complete configuration examples:

- [`web-server.sscfg`](examples/web-server.sscfg) - Complete web server setup
- [`database-server.sscfg`](examples/database-server.sscfg) - Database server configuration  
- [`docker-host.sscfg`](examples/docker-host.sscfg) - Docker host setup
- [`proxy-server.sscfg`](examples/proxy-server.sscfg) - Reverse proxy configuration
- [`build-server.sscfg`](examples/build-server.sscfg) - CI/CD build server setup

## ⚠️ Important Security Notes

1. **Always test on a VM first** - SetupSuite makes significant system changes
2. **Use strong SSH keys** - Never use the placeholder keys in production
3. **Review firewall rules** - Ensure only necessary ports are open
4. **Backup existing configs** - SetupSuite will backup original SSH configs
5. **Change default passwords** - Replace all placeholder passwords

## 🐳 Docker Log Fix

SetupSuite automatically implements the Docker logging fix mentioned in [`TODO.md`](TODO.md):

```json
{
  "log-driver": "local",
  "log-opts": {
    "max-size": "20m",
    "max-file": "5"
  }
}
```

This prevents Docker logs from consuming all disk space on busy systems.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 References

- [Crafting Interpreters](https://craftinginterpreters.com/contents.html) - Book about building code interpreters
- [Automated Server Setup Examples](https://github.com/ioniq-io/automated-server-setup)
- [Bash Aliases Collection](https://gist.githubusercontent.com/kemelzaidan/4e16d95e81db1ed90a4a/raw/1fdea821501265a9275ef36cfdc2a76d09cda5e5/.bash_aliases)

---

**⚡ Built for rapid server deployment and configuration management.**